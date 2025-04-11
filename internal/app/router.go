package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/krackl1n/golang-project/internal/handler"
	"github.com/krackl1n/golang-project/internal/middleware"
)

func getRouter(handler handler.Handler) *fiber.App {
	app := fiber.New()

	app.Use(middleware.MetricsMiddleware)

	userRouter := app.Group("/user")
	userRouter.Post("/", handler.CreateUser)
	userRouter.Get("/:id", handler.GetUser)
	userRouter.Delete("/:id", handler.DeleteUser)
	userRouter.Put("/", handler.UpdateUser)

	return app
}
