package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/krackl1n/golang-project/internal/metrics"
)

func MetricsMiddleware(c fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	duration := time.Since(start).Seconds()
	method := c.Method()
	endpoint := c.Path()
	statusCode := strconv.Itoa(c.Response().StatusCode())

	metrics.HttpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	metrics.HttpRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration)

	return err
}
