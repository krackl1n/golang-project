package cache

import (
	"context"
	"sync"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/metrics"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/krackl1n/golang-project/internal/repository"
)

type CacheDecorator struct {
	userRepository repository.UserProvider

	mu       sync.RWMutex
	users    map[uuid.UUID]models.User
	ttls     map[uuid.UUID]time.Time
	ttl      time.Duration
	stopChan chan struct{}
}

func New(userRepository repository.UserProvider, ttl time.Duration) *CacheDecorator {
	cache := &CacheDecorator{
		userRepository: userRepository,
		users:          make(map[uuid.UUID]models.User),
		ttls:           make(map[uuid.UUID]time.Time),
		ttl:            ttl,
		stopChan:       make(chan struct{}),
	}

	go cache.startWorker()

	return cache
}

func (c *CacheDecorator) Create(ctx context.Context, user *models.User) (uuid.UUID, error) {
	id, err := c.userRepository.Create(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	c.mu.Lock()
	c.users[id] = *user
	c.ttls[id] = time.Now().Add(c.ttl)
	metrics.CacheSize.Set(float64(unsafe.Sizeof(c.users)))
	c.mu.Unlock()

	return id, nil
}

func (c *CacheDecorator) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	c.mu.RLock()
	user, exists := c.users[id]
	expirationTime, ttlExists := c.ttls[id]
	c.mu.RUnlock()

	if exists && ttlExists && time.Now().Before(expirationTime) {
		metrics.CacheHits.Inc()
		return &user, nil
	}

	metrics.CacheMisses.Inc()
	userFromRepo, err := c.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.users[id] = *userFromRepo
	c.ttls[id] = time.Now().Add(c.ttl)
	metrics.CacheSize.Set(float64(unsafe.Sizeof(c.users)))
	c.mu.Unlock()

	return userFromRepo, nil
}

func (c *CacheDecorator) Update(ctx context.Context, user *models.User) error {
	err := c.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if _, exists := c.users[user.ID]; exists {
		c.users[user.ID] = *user
		c.ttls[user.ID] = time.Now().Add(c.ttl)
	}
	c.mu.Unlock()

	return nil
}

func (c *CacheDecorator) Delete(ctx context.Context, id uuid.UUID) error {
	err := c.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	c.mu.Lock()
	delete(c.users, id)
	delete(c.ttls, id)
	metrics.CacheSize.Set(float64(unsafe.Sizeof(c.users)))
	c.mu.Unlock()

	return nil
}

func (c *CacheDecorator) Stop() {
	close(c.stopChan)
}

func (c *CacheDecorator) startWorker() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredEntries()
		case <-c.stopChan:
			return
		}
	}
}

func (c *CacheDecorator) cleanupExpiredEntries() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, expirationTime := range c.ttls {
		if now.After(expirationTime) {
			delete(c.users, id)
			delete(c.ttls, id)
			metrics.CacheSize.Set(float64(unsafe.Sizeof(c.users)))
		}
	}
}
