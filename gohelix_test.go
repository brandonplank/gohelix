package gohelix

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestNew(t *testing.T) {
	new := New(&Options{ClientId: os.Getenv("CLIENT_ID"), ClientSecret: os.Getenv("CLIENT_SECRET")})
	_ = new.GetOAuthToken()
	log.Println(new.GetStream("lucroan").Data[0])
}
