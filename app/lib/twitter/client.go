package twitter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/mrjones/oauth"
)

func CreateTwitterClient() (*http.Client, error) {
	c := oauth.NewConsumer(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})
	c.Debug(true)

	t := oauth.AccessToken{
		Token:  os.Getenv("ACCESS_TOKEN_KEY"),
		Secret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}

	client, err := c.MakeHttpClient(&t)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateCRCToken(crcToken string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
