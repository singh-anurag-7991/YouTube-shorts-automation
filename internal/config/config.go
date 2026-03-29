package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PixabayAPIKey   string
	GoogleProjectID string
	WatermarkText   string
	YouTubePrivacy  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		PixabayAPIKey:   os.Getenv("PIXABAY_API_KEY"),
		GoogleProjectID: os.Getenv("GOOGLE_PROJECT_ID"),
		WatermarkText:   os.Getenv("CHANNEL_WATERMARK"),
		YouTubePrivacy:  os.Getenv("YOUTUBE_PRIVACY_STATUS"),
	}
}
