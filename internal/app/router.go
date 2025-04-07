package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krackl1n/golang-project/internal/handler"
)

func GetRouter(handler *handler.Handler) *fiber.App {
	app := fiber.New()

	app.Post("/user", handler.CreateUser)
	app.Get("/user/:id", handler.GetUser)
	app.Delete("/user/:id", handler.DeleteUser)
	app.Put("/user", handler.UpdateUser)

	return app
}
