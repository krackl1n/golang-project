package handler

import (
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
	createUserDTO := new(models.CreateUserDTO)
	if err := c.Bind().Body(createUserDTO); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body"),
		})
	}

	validate := validator.New()
	if err := validate.Struct(createUserDTO); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body"),
		})
	}

	idUser, err := h.userUseCase.CreateUser(createUserDTO)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed create user"),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id": idUser,
	})
}

func (h *Handler) GetUser(c fiber.Ctx) error {
	uuidUser, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed create user"),
		})
	}

	user, err := h.userUseCase.GetUserById(uuidUser)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": errors.Wrap(err, "user not found"),
		})
	}

	return c.Status(200).JSON(user)
}

func (h *Handler) UpdateUser(c fiber.Ctx) error {
	user := new(models.User)
	if err := c.Bind().Body(user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body"),
		})
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		errors.Wrap(err, "invalid request body")
	}

	if err := h.userUseCase.UpdateUser(user); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": errors.Wrap(err, "User not found"),
		})
	}

	return c.SendStatus(200)
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	uuidUser, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "failed create user"),
		})
	}

	if err = h.userUseCase.DeleteUser(uuidUser); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": errors.Wrap(err, "user not found"),
		})
	}

	return c.SendStatus(200)
}
