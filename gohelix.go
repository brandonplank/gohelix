package gohelix

import (
	"encoding/json"
	"errors"
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
	body, _, err := h.Request("POST", "https://id.twitch.tv/oauth2/token", payload, map[string]string{"Content-Type": "application/json"})
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
	h.setTokenIfNotValid()
	body, _, err := h.Request("GET", "https://api.twitch.tv/helix/streams?user_login="+name, nil, map[string]string{
		"Client-ID":     h.ClientId,
		"Authorization": "Bearer " + h.ClientOAuth,
	})
	if err != nil {
		log.Println(err.Error())
	}
	var stream Stream
	err = json.Unmarshal(body, &stream)
	if err != nil {
		log.Println(err.Error())
	}
	return stream
}

func (h *Helix) IsTokenValid() bool {
	_, status, err := h.Request("GET", "https://id.twitch.tv/oauth2/validate", nil, map[string]string{
		"Authorization": "OAuth " + h.ClientOAuth,
	})
	if err != nil {
		log.Println(err.Error())
	}
	if status == 401 && status != 200 {
		return false
	} else {
		return true
	}
}

func (h *Helix) setTokenIfNotValid() {
	if !h.IsTokenValid() {
		h.GetOAuthToken()
	}
}

func New(options *Options) (*Helix, error) {
	if len(options.ClientId) < 1 || len(options.ClientSecret) < 1 {
		return nil, errors.New("please set client_id and client_secret")
	}
	return &Helix{
		ClientId:     options.ClientId,
		ClientSecret: options.ClientSecret,
	}, nil
}
