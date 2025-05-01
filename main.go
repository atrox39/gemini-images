package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var (
	prompFlag string
)

func main() {
	flag.StringVar(&prompFlag, "promp", "A simple tortoise image", "Use promp to generate images")
	flag.Parse()
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	var imagesPath string = path.Join(currentPath, "images")
	if _, err := os.Stat(imagesPath); os.IsNotExist(err) {
		err := os.Mkdir(imagesPath, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	ctx := context.Background()
	godotenv.Load(".env")
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	parts := []*genai.Part{
		genai.NewPartFromText(prompFlag),
	}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp-image-generation",
		contents,
		config,
	)

	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			fmt.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			t := time.Now()
			outputFilename := fmt.Sprintf("%d%02d%02d_%02d%02d%02d.png", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
			_ = os.WriteFile(path.Join(imagesPath, outputFilename), imageBytes, 0644)
		}
	}
}
