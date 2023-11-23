package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/exec"
	"strconv"
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
	pubsub := redisClient.Subscribe(context.Background(), "hls")

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}

		log.Printf("[New message in channel (%s)]: %s", msg.Channel, msg.Payload)

		if msg.Channel == "hls" {
			inputFile := fmt.Sprintf("/videos/raw/%s", msg.Payload)
			outputDir := fmt.Sprintf("/videos/hls/%s", msg.Payload)

			go CreateHLS(inputFile, outputDir, 5)
		}
	}
}

func CreateHLS(inputFile string, outputDir string, segmentDuration int) {
	err := os.MkdirAll(outputDir, 0755)

	if err != nil {
		panic(fmt.Errorf("failed to create output directory: %v", err))
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
		panic(fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, output))
	}
}
