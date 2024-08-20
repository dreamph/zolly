package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c *fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, Test")
	})

	// Start the server on port 3001
	app.Listen(":3001")
}
