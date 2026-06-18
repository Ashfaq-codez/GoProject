package routes

import (
	"user-api/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	api := app.Group("/users")

	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)
}