package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"io"
	"io/fs"
	"log"
	"os"
)

var redisClient *redis.Client

func init() {
	rdC := redis.NewClient(&redis.Options{
		Addr:     "queue:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdC.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis connection failed. ", err)
	}

	log.Println("Redis successfully connected. Ping:", pong)

	redisClient = rdC
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit:         4 * 1024 * 1024,
		StreamRequestBody: true,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	app.Static("/videos", "/videos/hls")

	app.Post("/upload", func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			return err
		}

		titles, ok := form.Value["title"]
		if !ok {
			return errors.New("no title provided")
		}

		title := titles[0]
		if len(title) == 0 {
			return errors.New("title is empty")
		}

		videos, ok := form.File["video"]
		if !ok {
			return errors.New("no video provided")
		}

		if pathExists(fmt.Sprintf("/videos/hls/%s", title)) {
			return errors.New("a video with this title already exists")
		}

		video, err := videos[0].Open()
		if err != nil {
			return err
		}

		err = saveVideo(title, video)
		if err != nil {
			return err
		}

		redisClient.Publish(ctx.Context(), "hls", title)

		return nil
	})

	log.Fatal(app.Listen(":80"))
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}

func saveVideo(filename string, video io.Reader) error {
	path := fmt.Sprintf("/videos/raw/%s", filename)

	file, err := os.Create(path)
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}

	_, err = io.Copy(file, video)
	if err != nil {
		return err
	}

	return nil
}
