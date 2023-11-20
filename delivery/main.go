package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit:         4 * 1024 * 1024,
		StreamRequestBody: true,
	})

	app.Static("/", "./static")
	app.Static("/videos", "./hls")
	log.Fatal(app.Listen(":3000"))
}
