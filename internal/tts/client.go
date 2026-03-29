package tts

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

type Client struct {
	ctx context.Context
}

func NewClient() *Client {
	return &Client{
		ctx: context.Background(),
	}
}

func (c *Client) Synthesize(text string, outputPath string) error {
	client, err := texttospeech.NewClient(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to create texttospeech client: %w", err)
	}
	defer client.Close()

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Wavenet-D",
			SsmlGender:   texttospeechpb.SsmlGender_MALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(c.ctx, req)
	if err != nil {
		return fmt.Errorf("failed to synthesize speech: %w", err)
	}

	err = ioutil.WriteFile(outputPath, resp.AudioContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write audio file: %w", err)
	}

	log.Printf("✅ Audio content written to file: %v\n", outputPath)
	return nil
}
