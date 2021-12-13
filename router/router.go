package router

import (
	"ToDoApp/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	// index page
	index := router.Group("")
	index.Get("", handler.Index)

	// api routes
	api := router.Group("/api/v1")
	// CREATE
	api.Post("/create_todo", handler.NewRecord)

	// READ
	api.Get("/get_todos", handler.GetRecords)

	// UPDATE
	api.Put("/update_todo/:id", handler.UpdateRecord)

	// DELETE
	api.Delete("/delete_todo/:id", handler.DeleteRecord)

}
