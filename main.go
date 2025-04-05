package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

const (
	usersUrl = "/user"
	userUrl  = "/user/:id"
	port     = 8080
)

func main() {
	app := fiber.New()

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

// Можно еще сделать DTO для создания юзера
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
	repositoryUser *RepositotyUser
}

func NewHandler(repositoryUser *RepositotyUser) Handlers {
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

	// Возможно вынести валидацию в отдельную функцию
	name := c.Params("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing name",
		})
	}

	ageStr := c.Params("age")
	age, err := strconv.ParseInt(ageStr, 10, 8)
	if age > 200 || age <= 0 || err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid age",
		})
	}

	gender := c.Params("gender")
	//  Можно вынести отдельно
	validGenders := map[string]bool{"male": true, "female": true}
	if !validGenders[gender] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid gender",
		})
	}
	//  Не буду делать валидацию email:)
	email := c.Params("email")
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing email",
		})
	}

	// Желательно использовать отдельное DTO, но я не уверен в гошке делается это или нет
	idUser, err := (*h.repositoryUser).Create(User{
		Id:     0,
		Name:   name,
		Age:    uint8(age),
		Gender: gender,
		Email:  email,
	})

	if err != nil {
		// Я пока не знаю что отправить
		c.SendStatus(500)
	}

	return c.Status(201).JSON(fiber.Map{
		"id": int(idUser),
	})
}

func (h *handler) GetUser(c fiber.Ctx) error {
	idUserStr := c.Params("id")
	if idUserStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing id",
		})
	}

	idUser, err := strconv.ParseUint(idUserStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid id",
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
	if idUserStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing id",
		})
	}

	idUser, err := strconv.ParseUint(idUserStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid id",
		})
	}

	name := c.Params("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing name",
		})
	}

	ageStr := c.Params("age")
	age, err := strconv.ParseInt(ageStr, 10, 8)
	if age > 200 || age <= 0 || err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid age",
		})
	}

	gender := c.Params("gender")
	//  Можно вынести отдельно
	validGenders := map[string]bool{"male": true, "female": true}
	if !validGenders[gender] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid gender",
		})
	}
	//  Не буду делать валидацию email:)
	email := c.Params("email")
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing email",
		})
	}

	err = (*h.repositoryUser).Update(User{
		Id:     idUser,
		Name:   name,
		Age:    uint8(age),
		Gender: gender,
		Email:  email,
	})

	if err != nil {
		// TODO Расписать более подробно
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func (h *handler) DeleteUser(c fiber.Ctx) error {
	idUserStr := c.Params("id")
	if idUserStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing id",
		})
	}

	idUser, err := strconv.ParseUint(idUserStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid id",
		})
	}

	err = (*h.repositoryUser).Delete(idUser)
	if err != nil {
		// TODO Расписать более подробно
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

// repo

// Нужно ли передавать в репозитрий указатель на юзера или передавать копию
type RepositotyUser interface {
	Create(user User) (uint64, error)
	GetByID(id uint64) (*User, error)
	Update(user User) error
	Delete(id uint64) error
}

type repositotyUserLocal struct {
	//users map[int64]User
}

func NewRepositotyUserLocal() RepositotyUser {
	return &repositotyUserLocal{}
}

func (r *repositotyUserLocal) Create(user User) (uint64, error) {
	panic("unimplemented")
}

func (r *repositotyUserLocal) Delete(id uint64) error {
	panic("unimplemented")
}

func (r *repositotyUserLocal) GetByID(id uint64) (*User, error) {
	panic("unimplemented")
}

func (r *repositotyUserLocal) Update(user User) error {
	panic("unimplemented")
}
