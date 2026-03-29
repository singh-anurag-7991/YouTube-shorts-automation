package main

import (
	"fmt"
	"log"
	"youtube-shorts-automation/internal/config"
	"youtube-shorts-automation/internal/script"
)

func main() {
	fmt.Println("🚀 YouTube Shorts Automation - Processing Scripts")

	cfg := config.LoadConfig()
	if cfg.PixabayAPIKey == "" {
		log.Println("⚠️  Warning: PIXABAY_API_KEY is not set in .env")
	}

	mgr := script.NewManager("scripts.json")
	if err := mgr.Load(); err != nil {
		log.Fatalf("❌ Failed to load scripts: %v", err)
	}

	next, err := mgr.GetNext()
	if err != nil {
		log.Fatalf("❌ No available scripts: %v", err)
	}

	fmt.Printf("📝 Selected Script [%d]: %s\n", next.ID, next.Text)

	// Simulate marking as used
	if err := mgr.MarkAsUsed(next.ID); err != nil {
		log.Fatalf("❌ Failed to mark script as used: %v", err)
	}

	fmt.Println("✅ Script marked as used successfully")
}
