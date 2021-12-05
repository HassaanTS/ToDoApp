package main

import (
	"fmt"

	"ToDoApp/config"
	"ToDoApp/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	var err error
	// load environment
	err = config.LoadEnv()
	if err != nil {
		fmt.Println("error while loading environment...%w", err)
	}

	// starting fiber app
	app := fiber.New()

	// cors for resource sharing
	app.Use(cors.New())

	// set up router
	accessPoint := app.Group("/", logger.New())
	router.SetupRoutes(accessPoint)

	// serve app
	adress := fmt.Sprintf(":%s", config.GlobalConfig.AppPort)
	err = app.Listen(adress)
	if err != nil {
		fmt.Print("could not start server on ", adress)
	}
}
