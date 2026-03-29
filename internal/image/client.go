package image

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type PixabayResponse struct {
	Hits []struct {
		LargeImageURL string `json:"largeImageURL"`
		WebformatURL  string `json:"webformatURL"`
	} `json:"hits"`
}

type Client struct {
	APIKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
	}
}

func (c *Client) SearchAndDownload(query string, outputPath string) (string, error) {
	url := fmt.Sprintf("https://pixabay.com/api/?key=%s&q=%s&image_type=photo&orientation=vertical&per_page=20", c.APIKey, query)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call Pixabay API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Pixabay API returned status: %s", resp.Status)
	}

	var pixResp PixabayResponse
	if err := json.NewDecoder(resp.Body).Decode(&pixResp); err != nil {
		return "", fmt.Errorf("failed to decode Pixabay response: %w", err)
	}

	if len(pixResp.Hits) == 0 {
		return "", fmt.Errorf("no images found for query: %s", query)
	}

	// Pick a random hit
	rand.Seed(time.Now().UnixNano())
	hit := pixResp.Hits[rand.Intn(len(pixResp.Hits))]
	imageURL := hit.LargeImageURL
	if imageURL == "" {
		imageURL = hit.WebformatURL
	}

	// Download the image
	imgResp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer imgResp.Body.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, imgResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image content: %w", err)
	}

	return outputPath, nil
}
