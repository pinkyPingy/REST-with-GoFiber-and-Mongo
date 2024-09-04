package routes

import (
	"rest-with-gofiber-and-mongo/controllers"

	"github.com/gofiber/fiber/v2"
)

func EventRoute(app *fiber.App) {
	app.Post("/event", controllers.CreateEvent)
	app.Get("/event/:eventId", controllers.GetAnEvent)
	app.Put("/event/:eventId", controllers.EditAnEvent)
	app.Delete("/event/:eventId", controllers.DeleteAnEvent)
	app.Get("/events", controllers.GetAllEvents)
}
