package main

import (
	"fmt"
	"os"

	"ToDoApp/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("couldn't load environment from .env file...%w", err)
	}
	return nil
}

func main() {
	var err error
	// load environment
	err = loadEnv()
	if err != nil {
		fmt.Println("error while loading environment...%w", err)
	}

	// starting fiber app
	app := fiber.New()
	// set up router
	accessPoint := app.Group("/", logger.New())
	router.SetupRoutes(accessPoint)

	// serve app
	adress := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	err = app.Listen(adress)
	if err != nil {
		fmt.Print("could not start server on ", adress)
	}
}
