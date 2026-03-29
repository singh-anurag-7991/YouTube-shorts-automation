package main

import (
	"fmt"
	"log"
	"os"
	"youtube-shorts-automation/internal/config"
	"youtube-shorts-automation/internal/image"
	"youtube-shorts-automation/internal/script"
	"youtube-shorts-automation/internal/tts"
)

func main() {
	fmt.Println("🚀 YouTube Shorts Automation - Image & Audio Preparation")

	cfg := config.LoadConfig()
	if cfg.PixabayAPIKey == "" {
		log.Println("⚠️  Warning: PIXABAY_API_KEY is not set in .env")
	}

	// Ensure temp directory exists
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0755)
	}

	// 1. Script Selection
	mgr := script.NewManager("scripts.json")
	if err := mgr.Load(); err != nil {
		log.Fatalf("❌ Failed to load scripts: %v", err)
	}

	next, err := mgr.GetNext()
	if err != nil {
		log.Fatalf("❌ No available scripts: %v", err)
	}

	fmt.Printf("📝 Selected Script [%d]: %s\n", next.ID, next.Text)

	// 2. TTS Generation
	ttsClient := tts.NewClient()
	audioPath := fmt.Sprintf("temp/audio_%d.mp3", next.ID)

	log.Printf("⏳ Generating audio for script %d...", next.ID)
	if err := ttsClient.Synthesize(next.Text, audioPath); err != nil {
		log.Printf("❌ TTS failed: %v", err)
	}

	// 3. Image Download
	imgClient := image.NewClient(cfg.PixabayAPIKey)
	imagePath := fmt.Sprintf("temp/image_%d.jpg", next.ID)

	log.Printf("⏳ Fetching image for query 'god'...")
	downloadedPath, err := imgClient.SearchAndDownload("god", imagePath)
	if err != nil {
		log.Printf("❌ Image download failed: %v", err)
	} else {
		fmt.Printf("✅ Image downloaded: %s\n", downloadedPath)
	}

	// Simulate marking as used
	if err := mgr.MarkAsUsed(next.ID); err != nil {
		log.Fatalf("❌ Failed to mark script as used: %v", err)
	}

	fmt.Println("✅ Progress: Script marked as used, files in temp/")
}
