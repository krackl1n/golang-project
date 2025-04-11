package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/krackl1n/golang-project/internal/apperr"
	"github.com/krackl1n/golang-project/internal/models"
	"github.com/krackl1n/golang-project/internal/usecase"
	"github.com/pkg/errors"
)

type Handle struct {
	userUC usecase.UserProvider
}

func New(userUsecase usecase.UserProvider) Handler {
	return &Handle{
		userUC: userUsecase,
	}
}

func (h *Handle) CreateUser(c fiber.Ctx) error {
	createUserDTO := models.CreateUserDTO{}
	if err := c.Bind().Body(&createUserDTO); err != nil {
		slog.Debug("invalid request body", slog.Any("error", err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(createUserDTO); err != nil {
		slog.Debug("validate createUserDTO", slog.Any("error", err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	idUser, err := h.userUC.CreateUser(c.Context(), &createUserDTO)
	if err != nil {
		slog.Error("create user", slog.Any("error", err))
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id": idUser,
	})
}

func (h *Handle) GetUser(c fiber.Ctx) error {

	uuidUser, err := uuid.Parse(c.Params("id"))
	if err != nil {
		slog.Debug("parse id", slog.Any("error", err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "parse id").Error(),
		})
	}

	user, err := h.userUC.GetUserById(c.Context(), uuidUser)
	if err != nil {
		if errors.Is(err, apperr.ErrorNotFound) {
			slog.Debug("get by id", slog.Any("error", err))
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		slog.Error("get user", slog.Any("error", err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Wrap(err, "get user").Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(user)
}

func (h *Handle) UpdateUser(c fiber.Ctx) error {
	user := models.User{}
	if err := c.Bind().Body(&user); err != nil {
		slog.Debug("invalid request body", slog.Any("error", err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		slog.Debug("validate user", slog.Any("error", err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "invalid request body").Error(),
		})
	}

	if err := h.userUC.UpdateUser(c.Context(), &user); err != nil {
		if errors.Is(err, apperr.ErrorNotFound) {
			slog.Debug("update user", slog.Any("error", err))
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		slog.Error("update user", slog.Any("error", err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Wrap(err, "update user").Error(),
		})
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handle) DeleteUser(c fiber.Ctx) error {
	uuidUser, err := uuid.Parse(c.Params("id"))
	if err != nil {
		slog.Debug("delete user", slog.Any("error", err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errors.Wrap(err, "parse id").Error(),
		})
	}

	if err = h.userUC.DeleteUser(c.Context(), uuidUser); err != nil {
		if errors.Is(err, apperr.ErrorNotFound) {
			slog.Debug("update user", slog.Any("error", err))
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		slog.Error("delete user", slog.Any("error", err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": errors.Wrap(err, "delete user").Error(),
		})
	}

	return c.SendStatus(http.StatusOK)
}
