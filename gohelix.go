package gohelix

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Helix struct {
	ClientId     string
	ClientSecret string
	ClientOAuth  string

	HttpClient http.Client
}

type Options struct {
	ClientId     string
	ClientSecret string
}

type OAuthPayload struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Stream struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		UserLogin    string    `json:"user_login"`
		UserName     string    `json:"user_name"`
		GameID       string    `json:"game_id"`
		GameName     string    `json:"game_name"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailURL string    `json:"thumbnail_url"`
		TagIds       []string  `json:"tag_ids"`
		IsMature     bool      `json:"is_mature"`
	} `json:"data"`
	Pagination struct {
	} `json:"pagination"`
}

func (h *Helix) GetOAuthToken() Token {
	h.HttpClient = http.Client{}
	payload, _ := json.MarshalIndent(OAuthPayload{
		ClientId:     h.ClientId,
		ClientSecret: h.ClientSecret,
		GrantType:    "client_credentials",
	}, "", "\t")
	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}
	res, err := h.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println(err.Error())
	}
	h.ClientOAuth = token.AccessToken
	return token
}

func (h *Helix) GetStream(name string) Stream {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_login="+name, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header = http.Header{
		"Client-ID":     []string{"m2pipnmdfk93o5pexx1h19vowkpkxp"},
		"Authorization": []string{"Bearer " + h.ClientOAuth},
	}

	res, err := h.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var stream Stream
	err = json.Unmarshal(body, &stream)
	if err != nil {
		log.Println(err.Error())
	}
	return stream
}

func New(options *Options) *Helix {
	return &Helix{
		ClientId:     options.ClientId,
		ClientSecret: options.ClientSecret,
	}
}
