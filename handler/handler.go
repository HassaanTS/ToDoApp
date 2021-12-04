package handler

import (
	"ToDoApp/config"
	"ToDoApp/db"
	"ToDoApp/todos"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Index(c *fiber.Ctx) error {
	err := c.SendString("Fiber is alive!")
	return err
}

func GetRecords(c *fiber.Ctx) error {
	// open and close connection so the handle doesn't get stale
	client, ctx, cancel, err := db.ConnectDB(db.BuildURI())
	if err != nil {
		return fmt.Errorf("couldn't connect to database")
	}
	defer db.DisconnectDB(client, ctx, cancel)

	records, err := db.GetRecords(client, ctx, config.GlobalConfig.MongoDBName, config.GlobalConfig.MongoDBCollection)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": fmt.Errorf("error while fetching records from the database...%w", err).Error(),
		}))
	}

	// return id of record to client
	return c.JSON((fiber.Map{
		"records": records,
	}))
}

func NewRecord(c *fiber.Ctx) error {
	// open and close connection so the handle doesn't get stale
	client, ctx, cancel, err := db.ConnectDB(db.BuildURI())
	if err != nil {
		return fmt.Errorf("couldn't connect to database")
	}
	defer db.DisconnectDB(client, ctx, cancel)

	// create record
	t := todos.New()

	// parse body and populate ToDo
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": fmt.Errorf("can't parse request body while creating new record...%w", err),
		}))
	}

	// insert record in db
	insertOneResult, err := db.InsertRecord(client, ctx, config.GlobalConfig.MongoDBName, config.GlobalConfig.MongoDBCollection, t)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": fmt.Errorf("error while inserting record in the database...%w", err),
		}))
	}

	// return id of record to client
	return c.JSON((fiber.Map{
		"new_record_id": insertOneResult.InsertedID,
	}))
}

func UpdateRecord(c *fiber.Ctx) error {
	// open and close connection so the handle doesn't get stale
	client, ctx, cancel, err := db.ConnectDB(db.BuildURI())
	if err != nil {
		return fmt.Errorf("couldn't connect to database")
	}
	defer db.DisconnectDB(client, ctx, cancel)

	// id needs to be primitive id type to actually filter db
	id := c.Params("id")
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	// parse body for changes
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": fmt.Errorf("can't parse request body while updating record...%w", err),
		}))
	}

	// convert changes to bson document
	var changes bson.D
	for k, v := range body {
		changes = bson.D{{k, v}}
	}
	update := bson.D{{"$set", changes}}

	// update the record
	result, err := db.UpdateRecord(client, ctx, config.GlobalConfig.MongoDBName, config.GlobalConfig.MongoDBCollection, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": fmt.Errorf("error while updating record...%w", err),
		}))
	}

	// return changes
	return c.JSON((fiber.Map{
		"records_changed": result.ModifiedCount,
	}))
}

func DeleteRecord(c *fiber.Ctx) error {
	// open and close connection so the handle doesn't get stale
	client, ctx, cancel, err := db.ConnectDB(db.BuildURI())
	if err != nil {
		return fmt.Errorf("couldn't connect to database")
	}
	defer db.DisconnectDB(client, ctx, cancel)

	// id needs to be primitive id type to actually filter db
	id := c.Params("id")
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	// update the record
	result, err := db.DeleteRecord(client, ctx, config.GlobalConfig.MongoDBName, config.GlobalConfig.MongoDBCollection, filter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON((fiber.Map{
			"error": fmt.Errorf("error while deleting record...%w", err),
		}))
	}

	// return changes
	return c.JSON((fiber.Map{
		"records_deleted": result.DeletedCount,
	}))
}
