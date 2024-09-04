package main

import (
	"log"
	"rest-with-gofiber-and-mongo/configs"
	"rest-with-gofiber-and-mongo/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.EventRoute(app)

	log.Fatal(app.Listen(":6000"))
}
