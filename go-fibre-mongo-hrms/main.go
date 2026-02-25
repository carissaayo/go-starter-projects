package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017" + dbName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id, omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return err
	}

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     client.Database(dbName),
	}

	log.Println("✅ Connected to MongoDB")
	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := mg.Db.Collection("employees")

		filter := bson.D{}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var employees []Employee

		if err := cursor.All(ctx, &employees); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(employees)
	})

	app.Post("/employee", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := mg.Db.Collection("employees")

		var employee Employee

		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		employee.ID = primitive.NewObjectID().Hex()

		result, err := collection.InsertOne(ctx, employee)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(201).JSON(result)

	})

	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := mg.Db.Collection("employees")

		idParam := c.Params("id")

		employeeID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid employee ID",
			})
		}

		var employee Employee

		if err := c.BodyParser(&employee); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		filter := bson.M{"_id": employeeID}

		update := bson.M{
			"$set": bson.M{
				"name":   employee.Name,
				"age":    employee.Age,
				"salary": employee.Salary,
			},
		}

		opts := options.FindOneAndUpdate().
			SetReturnDocument(options.After)

		var updatedEmployee Employee

		err = collection.FindOneAndUpdate(ctx, filter, update, opts).
			Decode(&updatedEmployee)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).JSON(fiber.Map{
					"error": "Employee not found",
				})
			}

			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(updatedEmployee)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		collection := mg.Db.Collection("employees")

		idParam := c.Params("id")

		employeeID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid employee ID",
			})
		}

		filter := bson.M{"_id": employeeID}

		result, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if result.DeletedCount == 0 {
			return c.Status(404).JSON(fiber.Map{
				"error": "Employee not found",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Employee deleted successfully",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
