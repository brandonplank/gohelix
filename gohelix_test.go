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
	new, _ := New(&Options{ClientId: os.Getenv("CLIENT_ID"), ClientSecret: os.Getenv("CLIENT_SECRET")})
	_ = new.GetOAuthToken()
	log.Println(new.GetStream("lifestomper"))
}

func TestHelix_GetOAuthToken(t *testing.T) {
	new, _ := New(&Options{ClientId: os.Getenv("CLIENT_ID"), ClientSecret: os.Getenv("CLIENT_SECRET")})
	log.Println(new.GetOAuthToken())
}
