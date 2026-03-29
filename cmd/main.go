package main

import (
	"fmt"
	"log"
	"youtube-shorts-automation/internal/config"
)

func main() {
	fmt.Println("🚀 YouTube Shorts Automation - Project Scaffold Initialized")

	cfg := config.LoadConfig()
	if cfg.PixabayAPIKey == "" {
		log.Println("⚠️  Warning: PIXABAY_API_KEY is not set in .env")
	}

	fmt.Println("✅ Configuration loaded successfully")
}
