package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit:         4 * 1024 * 1024,
		StreamRequestBody: true,
	})

	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("video")
		if err != nil {
			return err
		}

		fmt.Println(file.Header)

		src := fmt.Sprintf("./temp/%s", file.Filename)

		err = ctx.SaveFile(file, src)
		if err != nil {
			return err
		}

		return CreateHLS(src, "./hls/", 5)
	})

	app.Listen(":3001")
}

func CreateHLS(inputFile string, outputDir string, segmentDuration int) error {
	err := os.MkdirAll(outputDir, 0755)

	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-profile:v", "baseline",
		"-level", "3.0",
		"-start_number", "0", // start numbering segments from 0
		"-hls_time", strconv.Itoa(segmentDuration), // duration of each segment in seconds
		"-hls_list_size", "0", // keep all segments in the playlist
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", outputDir),
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, output)
	}

	return nil
}
