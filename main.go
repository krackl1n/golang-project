package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v3"
)

const (
	usersUrl = "/user"
	userUrl  = "/user/:id"
	port     = 8080
)

func main() {
	app := fiber.New()

	//  TODO Добавить логирование
	log.Println("init fiber app")
	start(app)
}

func start(app *fiber.App) {
	repositoryUser := NewRepositotyUserLocal()

	handler := NewHandler(&repositoryUser)
	handler.Register(app)

	log.Println("listening port ", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}

// entity

// Можно еще сделать DTO для создания юзера(даже нужно)
type User struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	Age    uint8  `json:"age"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
}

//  handlers

type Handlers interface {
	Register(app *fiber.App)
}

type handler struct {
	repositoryUser *RepositoryUser
}

func NewHandler(repositoryUser *RepositoryUser) Handlers {
	return &handler{
		repositoryUser: repositoryUser,
	}
}

func (h *handler) Register(app *fiber.App) {
	app.Post(usersUrl, h.CreateUser)
	app.Get(userUrl, h.GetUser)
	app.Delete(userUrl, h.DeleteUser)
	app.Put(userUrl, h.UpdateUser)
}

func (h *handler) CreateUser(c fiber.Ctx) error {
	name := c.Params("name")
	ageStr := c.Params("age")
	gender := c.Params("gender")
	//  Не буду делать валидацию email:)
	email := c.Params("email")

	age, err := validateUser(name, ageStr, gender, email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Желательно использовать отдельное DTO, но я не уверен в гошке делается это или нет
	idUser, err := (*h.repositoryUser).Create(&User{
		Name:   name,
		Age:    uint8(age),
		Gender: gender,
		Email:  email,
	})

	if err != nil {
		// Я пока не знаю что отправить, пусть будет на будущее)
		c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("server error: %v", err),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"id": int(idUser),
	})
}

func (h *handler) GetUser(c fiber.Ctx) error {
	idUserStr := c.Params("id")
	idUser, err := validateIdUser(idUserStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := (*h.repositoryUser).GetByID(idUser)

	if user == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(user)
}

func (h *handler) UpdateUser(c fiber.Ctx) error {
	idUserStr := c.Params("id")
	name := c.Params("name")
	ageStr := c.Params("age")
	gender := c.Params("gender")
	//  Не буду делать валидацию email:)
	email := c.Params("email")

	idUser, err := validateIdUser(idUserStr)
	age, err1 := validateUser(name, ageStr, gender, email)
	// Нормально ли так делать?
	if err != nil || err1 != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = (*h.repositoryUser).Update(&User{
		Id:     idUser,
		Name:   name,
		Age:    age,
		Gender: gender,
		Email:  email,
	})

	if err != nil {
		// TODO Расписать более подробно
		return c.Status(500).JSON(fiber.Map{
			"server error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func (h *handler) DeleteUser(c fiber.Ctx) error {
	idUserStr := c.Params("id")
	idUser, err := validateIdUser(idUserStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = (*h.repositoryUser).Delete(idUser)
	if err != nil {
		// TODO Расписать более подробно
		return c.Status(500).JSON(fiber.Map{
			"server error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func validateIdUser(idUserStr string) (uint64, error) {
	if idUserStr == "" {
		return 0, fmt.Errorf("missing id")
	}

	idUser, err := strconv.ParseUint(idUserStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}

	return idUser, nil
}

func validateUser(name string, ageStr string, gender string, email string) (uint8, error) {
	if name == "" {
		return 0, fmt.Errorf("missing name")
	}

	age, err := strconv.ParseUint(ageStr, 10, 8)
	if err != nil || age > 200 || age <= 0 {
		return 0, fmt.Errorf("invalid age")
	}

	validGenders := map[string]bool{"male": true, "female": true}
	if !validGenders[gender] {
		return 0, fmt.Errorf("invalid gender")
	}

	if email == "" {
		return 0, fmt.Errorf("missing email")
	}

	return uint8(age), nil
}

// repo

// Нужно ли передавать в репозитрий указатель на юзера или передавать копию?
type RepositoryUser interface {
	Create(user *User) (uint64, error)
	GetByID(id uint64) (*User, error)
	Update(user *User) error
	Delete(id uint64) error
}

type repositoryUserLocal struct {
	users  map[uint64]User
	nextID uint64
	mu     sync.Mutex
}

func NewRepositotyUserLocal() RepositoryUser {
	return &repositoryUserLocal{
		users:  make(map[uint64]User),
		nextID: 0,
	}
}

func (r *repositoryUserLocal) Create(user *User) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	idUser := r.nextID
	r.users[idUser] = *user
	r.nextID++

	return idUser, nil
}

func (r *repositoryUserLocal) GetByID(id uint64) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (r *repositoryUserLocal) Update(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[user.Id]
	if !exists {
		return fmt.Errorf("user not found")
	}

	r.users[user.Id] = *user
	return nil
}

func (r *repositoryUserLocal) Delete(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[id]
	if !exists {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)
	return nil
}
