package main

import (
	"fmt"
	"log"
	"os"
	"youtube-shorts-automation/internal/config"
	"youtube-shorts-automation/internal/script"
	"youtube-shorts-automation/internal/tts"
)

func main() {
	fmt.Println("🚀 YouTube Shorts Automation - Audio Generation")

	cfg := config.LoadConfig()
	if cfg.PixabayAPIKey == "" {
		log.Println("⚠️  Warning: PIXABAY_API_KEY is not set in .env")
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

	// Ensure temp directory exists
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0755)
	}

	log.Printf("⏳ Generating audio for script %d...", next.ID)
	if err := ttsClient.Synthesize(next.Text, audioPath); err != nil {
		log.Printf("❌ TTS failed (Check GOOGLE_APPLICATION_CREDENTIALS): %v", err)
	} else {
		fmt.Printf("✅ Audio generated: %s\n", audioPath)
	}

	// Simulate marking as used
	if err := mgr.MarkAsUsed(next.ID); err != nil {
		log.Fatalf("❌ Failed to mark script as used: %v", err)
	}

	fmt.Println("✅ Script marked as used successfully")
}
