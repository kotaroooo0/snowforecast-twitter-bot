package key

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func GetTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("API_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("API_SECRET_KEY"))
	api := anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	return api
}

func CreateCRCToken(crcToken string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("API_KEY")))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
