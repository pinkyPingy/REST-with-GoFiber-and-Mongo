package controllers

import (
	"context"
	"net/http"
	"rest-with-gofiber-and-mongo/configs"
	"rest-with-gofiber-and-mongo/models"
	"rest-with-gofiber-and-mongo/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventCollection *mongo.Collection = configs.GetCollection(configs.DB, "events")
var validate = validator.New()

func CreateEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var event models.Event
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&event); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&event); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newEvent := models.Event{
		// Id:          primitive.NewObjectID(),
		Title:       event.Title,
		Description: event.Description,
		Date:        event.Date,
		Time:        event.Time,
		Location:    event.Location,
		Amount:      event.Amount,
	}

	result, err := eventCollection.InsertOne(ctx, newEvent)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.EventResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAnEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	eventId := c.Params("eventId")
	var event models.Event
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	err := eventCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&event)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": event}})
}

func EditAnEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	eventId := c.Params("eventId")
	var event models.Event
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	//validate the request body
	if err := c.BodyParser(&event); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&event); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.EventResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"Title":       event.Title,
		"Description": event.Description,
		"Date":        event.Date,
		"Time":        event.Time,
		"Location":    event.Location,
		"Amount":      event.Amount,
	}

	result, err := eventCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated event details
	var updatedEvent models.Event
	if result.MatchedCount == 1 {
		err := eventCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedEvent)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedEvent}})
}

func DeleteAnEvent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	eventId := c.Params("eventId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(eventId)

	result, err := eventCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.EventResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Event with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Event successfully deleted!"}},
	)
}

func GetAllEvents(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var events []models.Event
	defer cancel()

	results, err := eventCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleEvent models.Event
		if err = results.Decode(&singleEvent); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.EventResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		events = append(events, singleEvent)
	}

	return c.Status(http.StatusOK).JSON(
		responses.EventResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": events}},
	)
}
