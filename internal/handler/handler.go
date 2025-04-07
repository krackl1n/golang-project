package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/krackl1n/golang-project/internal/usecase"
	"github.com/pkg/errors"
)

type Handler struct {
	userUseCase *usecase.UseCase
}

func NewHandler(userUseCase *usecase.UseCase) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

func (h *Handler) CreateUser(c fiber.Ctx) error {
	// TODO Логи переделать
	createUserDTO := new(models.CreateUserDTO)

	if err := c.Bind().Body(createUserDTO); err != nil {
		slog.Warn("Invalid request body", "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	validate := validator.New()
	if err := validate.Struct(createUserDTO); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	// Потом добавлю контекст
	ctx := context.TODO()
	idUser, err := h.userUseCase.CreateUser(ctx, createUserDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed create user").Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id": idUser,
	})
}

func (h *Handler) GetUser(c fiber.Ctx) error {
	// TODO Логи переделать
	uuidUser, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed create user").Error(),
		})
	}

	// Потом добавлю контекст
	ctx := context.TODO()
	user, err := h.userUseCase.GetUserById(ctx, uuidUser)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": errors.Wrap(err, "user not found").Error(),
		})
	}

	return c.Status(200).JSON(user)
}

func (h *Handler) UpdateUser(c fiber.Ctx) error {
	// TODO Логи переделать
	user := new(models.User)
	if err := c.Bind().Body(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	// Потом добавлю контекст
	ctx := context.TODO()
	if err := h.userUseCase.UpdateUser(ctx, user); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": errors.Wrap(err, "user not found").Error(),
		})
	}

	return c.SendStatus(200)
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	// TODO Логи переделать
	uuidUser, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed create user").Error(),
		})
	}

	// Потом добавлю контекст
	ctx := context.TODO()
	if err = h.userUseCase.DeleteUser(ctx, uuidUser); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": errors.Wrap(err, "user not found").Error(),
		})
	}

	return c.SendStatus(200)
}
