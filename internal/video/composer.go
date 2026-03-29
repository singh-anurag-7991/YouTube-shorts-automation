package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Composer struct {
	Watermark string
}

func NewComposer(watermark string) *Composer {
	return &Composer{
		Watermark: watermark,
	}
}

func (c *Composer) CreateShort(imagePath, audioPath, outputPath string) error {
	// Ensure input files exist
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file does not exist: %s", imagePath)
	}
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return fmt.Errorf("audio file does not exist: %s", audioPath)
	}

	// 1. Process image: Scale and crop to 1080x1920 (9:16)
	// 2. Loop image to match audio duration
	// 3. Overlay audio
	// 4. Add text watermark

	log.Printf("🎬 Composing video: %s + %s -> %s", imagePath, audioPath, outputPath)

	videoStream := ffmpeg.Input(imagePath, ffmpeg.KwArgs{"loop": 1}).
		Filter("scale", ffmpeg.Args{"1080", "1920", "force_original_aspect_ratio=increase"}).
		Filter("crop", ffmpeg.Args{"1080", "1920"}).
		Filter("drawtext", ffmpeg.Args{
			fmt.Sprintf("text='%s':fontsize=48:fontcolor=white@0.4:x=w-tw-50:y=h-th-50", c.Watermark),
		})

	audioStream := ffmpeg.Input(audioPath)

	err := ffmpeg.Output([]*ffmpeg.Stream{videoStream, audioStream}, outputPath, ffmpeg.KwArgs{
		"shortest": "",
		"pix_fmt":  "yuv420p",
		"vcodec":   "libx264",
		"acodec":   "aac",
	}).OverWriteOutput().Run()

	if err != nil {
		return fmt.Errorf("ffmpeg command failed: %w", err)
	}

	log.Printf("✅ Video composed successfully: %s", outputPath)
	return nil
}

// CheckFFmpeg looks for the ffmpeg binary in the system path.
func CheckFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}
