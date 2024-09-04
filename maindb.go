package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello World")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MONGO_URI := os.Getenv("DB_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	collection = client.Database("golang_todo").Collection("todos")
	app := fiber.New()

	app.Get("/api/todos", GetTodos)
	app.Post("/api/todos", PostTodo)
	app.Patch("/api/todos/:id", PatchTodo)
	app.Delete("/api/todos/:id", DeleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Fatal(app.Listen(":" + port))

}

func GetTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	return c.Status(200).JSON(todos)
}

func PostTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	err := c.BodyParser(todo)
	if err != nil {
		log.Fatal(err)
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"msg": "Body is empty"})
	}

	one, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = one.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)

}

func PatchTodo(c *fiber.Ctx) error {

	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": bson.M{"completed": true}})

	if err != nil {
		log.Fatal(err)
	}
	return c.Status(200).JSON(fiber.Map{"success": true})

}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		log.Fatal(err)
	}
	return c.Status(200).JSON(fiber.Map{"Success": true})
}
