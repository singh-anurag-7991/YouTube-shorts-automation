package main

import (
	"fmt"
	"os"
	"time"
	"youtube-shorts-automation/internal/config"
	"youtube-shorts-automation/internal/runner"
)

func main() {
	fmt.Println("🌟 YouTube Shorts Automation - Runner Active")

	cfg := config.LoadConfig()

	// Ensure temp directory exists
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0755)
	}

	r := runner.NewRunner(cfg)

	// User requirement: twice or thrice a day.
	// 24 hours / 3 = 8 hours.
	interval := 8 * time.Hour

	// For testing purposes, you might want to call r.RunOnce() directly.
	// We'll start the scheduler here.
	r.StartScheduler(interval)
}
