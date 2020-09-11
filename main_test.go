package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joho/godotenv"
)

func before() {
	err := godotenv.Load(".env.sample")
	if err != nil {
		log.Fatal("Error loading .env.sample file")
	}
}

func TestGetTwitterWebhook(t *testing.T) {
	before()

	router, err := setupRouter()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	values := url.Values{}
	values.Add("crc_token", "test")
	req, _ := http.NewRequest("GET", "twitter_webhook", nil)
	req.URL.RawQuery = values.Encode()
	router.ServeHTTP(w, req)

	byteArray, _ := ioutil.ReadAll(w.Body)

	data := TokenResponse{}
	if err := json.Unmarshal(([]byte)(byteArray), &data); err != nil {
		t.Errorf("JSON Unmarshal error: %s", err)
	}
	if diff := cmp.Diff(w.Code, 200); diff != "" {
		t.Errorf("Diff: (-got +want)\n%s", diff)
	}
	if diff := cmp.Diff(data.Token, "sha256=KDMuKWcx/Rw8ofrHYjc5atBXnxT4mjqjL9MfGvrY8j4="); diff != "" {
		t.Errorf("Diff: (-got +want)\n%s", diff)
	}
}

type TokenResponse struct {
	Token string `json:"response_token"`
}
