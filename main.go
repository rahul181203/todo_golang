package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

//type Todo struct {
//	ID        int    `json:"id" bson:"_id"`
//	Completed bool   `json:"completed"`
//	Body      string `json:"body"`
//}

func main() {
	fmt.Println("Hello World")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	todos := []Todo{}

	port := os.Getenv("PORT")

	app.Get("/api/todo", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	app.Post("/api/todo", func(c *fiber.Ctx) error {
		todo := &Todo{} // { id:0, completed:false, body:"" }
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "body is empty"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(201).JSON(todo)
	})

	app.Patch("/api/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "id not found"})
	})

	app.Delete("/api/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(todos)
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "id not found"})

	})

	log.Fatal(app.Listen(":" + port))
}
