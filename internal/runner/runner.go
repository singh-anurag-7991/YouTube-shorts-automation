package runner

import (
	"fmt"
	"log"
	"time"
	"youtube-shorts-automation/internal/config"
	"youtube-shorts-automation/internal/image"
	"youtube-shorts-automation/internal/script"
	"youtube-shorts-automation/internal/tts"
	"youtube-shorts-automation/internal/video"
	"youtube-shorts-automation/internal/youtube"
)

type Runner struct {
	Config    *config.Config
	ScriptMgr *script.Manager
	TTS       *tts.Client
	Image     *image.Client
	Video     *video.Composer
	YouTube   *youtube.Uploader
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		Config:    cfg,
		ScriptMgr: script.NewManager("scripts.json"),
		TTS:       tts.NewClient(),
		Image:     image.NewClient(cfg.PixabayAPIKey),
		Video:     video.NewComposer(cfg.WatermarkText),
		YouTube:   youtube.NewUploader(),
	}
}

func (r *Runner) RunOnce() error {
	log.Println("🚀 Starting automation cycle...")

	// 1. Load Scripts
	if err := r.ScriptMgr.Load(); err != nil {
		return fmt.Errorf("script load failed: %v", err)
	}

	// 2. Get Next Script
	next, err := r.ScriptMgr.GetNext()
	if err != nil {
		return fmt.Errorf("no scripts available: %v", err)
	}
	log.Printf("📝 Selected Script [%d]: %s", next.ID, next.Text)

	// 3. Audio & Image Paths
	audioPath := fmt.Sprintf("temp/audio_%d.mp3", next.ID)
	imagePath := fmt.Sprintf("temp/image_%d.jpg", next.ID)
	videoPath := fmt.Sprintf("temp/short_%d.mp4", next.ID)

	// 4. Generate Audio
	if err := r.TTS.Synthesize(next.Text, audioPath); err != nil {
		log.Printf("⚠️  TTS skipped: %v", err)
	}

	// 5. Download Image
	if _, err := r.Image.SearchAndDownload("god", imagePath); err != nil {
		log.Printf("⚠️  Image download skipped: %v", err)
	}

	// 6. Compose Video
	if err := r.Video.CreateShort(imagePath, audioPath, videoPath); err != nil {
		log.Printf("⚠️  Composition skipped: %v", err)
	}

	// 7. Upload to YouTube
	if err := r.YouTube.Init("client_secret.json"); err != nil {
		log.Printf("⚠️  YouTube Init skipped: %v", err)
	} else {
		title := fmt.Sprintf("God Message #%d", next.ID)
		videoID, err := r.YouTube.Upload(videoPath, title, next.Text, r.Config.YouTubePrivacy)
		if err != nil {
			log.Printf("❌ YouTube upload failed: %v", err)
		} else {
			log.Printf("✅ Uploaded! ID: %s", videoID)
		}
	}

	// 8. Mark as Used
	if err := r.ScriptMgr.MarkAsUsed(next.ID); err != nil {
		return fmt.Errorf("failed to mark script as used: %v", err)
	}

	log.Println("✅ Cycle completed successfully")
	return nil
}

func (r *Runner) StartScheduler(interval time.Duration) {
	log.Printf("⏰ Starting scheduler with interval: %v", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run once immediately
	if err := r.RunOnce(); err != nil {
		log.Printf("❌ Run error: %v", err)
	}

	for range ticker.C {
		if err := r.RunOnce(); err != nil {
			log.Printf("❌ Run error: %v", err)
		}
	}
}
