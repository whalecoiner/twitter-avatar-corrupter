package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
)

func getOrginalProfileURL(profileURL string) string {
	return strings.Replace(profileURL, "_normal.", ".", 1)
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	godotenv.Load()

	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	screenName := "whalecoiner"

	log.Println(consumerKey)
	log.Println(consumerSecret)
	log.Println(accessToken)
	log.Println(accessSecret)

	if consumerKey == "" {
		log.Fatal("TWITTER_CONSUMER_KEY must be set")
	}

	if consumerSecret == "" {
		log.Fatal("TWITTER_CONSUMER_SECRET must be set")
	}

	if accessToken == "" {
		log.Fatal("TWITTER_ACCESS_TOKEN must be set")
	}

	if accessSecret == "" {
		log.Fatal("TWITTER_ACCESS_SECRET must be set")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	client := twitter.NewClient(httpClient)

	params := &twitter.UserShowParams{ScreenName: screenName}
	user, _, err := client.Users.Show(params)

	if err != nil {
		log.Fatal(err)
	}

	originalProfileImageURL := getOrginalProfileURL(user.ProfileImageURL)
	filename := path.Base(originalProfileImageURL)

	log.Println("Profile Image URL", originalProfileImageURL)

	downloadErr := DownloadFile(filename, originalProfileImageURL)
	if err != nil {
		log.Fatal(downloadErr)
	}

	img, imageOpenErr := imaging.Open(filename)

	if imageOpenErr != nil {
		log.Fatal(imageOpenErr)
	}

	img = imaging.AdjustSaturation(img, -100)

	imageSaveErr := imaging.Save(img, filename)

	if imageSaveErr != nil {
		log.Fatal(imageSaveErr)
	}

}
