package app

import (
	"fmt"

	"github.com/krackl1n/golang-project/internal/handler"
	"github.com/krackl1n/golang-project/internal/repository"
	"github.com/krackl1n/golang-project/internal/usecase"
)

// Вынести в конфиг
const (
	connString = "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres"
)

func Run() error {
	//  TODO Добавить логирование
	// if err := database.Migrate(connString); err != nil {
	// 	return err
	// }

	// log.Println("init fiber app")

	userRepository := repository.NewUserRepository()
	usecase := usecase.NewUseCase(userRepository)
	handler := handler.NewHandler(usecase)

	app := GetRouter(handler)

	// log.Println("listening port ", 8080)

	err := app.Listen(fmt.Sprintf(":%d", 8080))
	return err
}
