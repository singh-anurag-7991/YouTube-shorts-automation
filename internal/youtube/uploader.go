package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Uploader struct {
	ctx     context.Context
	service *youtube.Service
}

func NewUploader() *Uploader {
	return &Uploader{
		ctx: context.Background(),
	}
}

func (u *Uploader) Init(clientSecretPath string) error {
	b, err := os.ReadFile(clientSecretPath)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	// Scope for YouTube upload
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := u.getClient(config)
	service, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve YouTube client: %v", err)
	}

	u.service = service
	return nil
}

func (u *Uploader) Upload(filename string, title string, description string, privacy string) (string, error) {
	if u.service == nil {
		return "", fmt.Errorf("youtube service not initialized")
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: description,
			CategoryId:  "22", // People & Blogs
		},
		Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	call := u.service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("error uploading video: %v", err)
	}

	return response.Id, nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func (u *Uploader) getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := u.tokenFromFile(tokenFile)
	if err != nil {
		tok = u.getTokenFromWeb(config)
		u.saveToken(tokenFile, tok)
	}
	return config.Client(u.ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
func (u *Uploader) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(u.ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile retrieves a Token from a given file path.
func (u *Uploader) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file path.
func (u *Uploader) saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
