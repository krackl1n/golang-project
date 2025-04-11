package handler

import "github.com/gofiber/fiber/v3"

type Handler interface {
	CreateUser(c fiber.Ctx) error
	GetUser(c fiber.Ctx) error
	UpdateUser(c fiber.Ctx) error
	DeleteUser(c fiber.Ctx) error
}
