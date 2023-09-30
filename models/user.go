package models

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

type User struct {
	DisplayName     string `json:"display_name"`
	BrodcasterLogin string `json:"broadcaster_login"`
	IsLive          bool   `json:"is_live"`
	ThumbnailUrl    string `json:"thumbnail_url"`
}

func Search(user_name string) (*User, error) {
	var u User
	url := "https://api.twitch.tv/helix/search/channels?query=" + user_name

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req.Header.Add("Client-Id", os.Getenv("CLIENT_ID"))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var users struct {
		Users []User `json:"data"`
	}

	json.Unmarshal(body, &users)

	for _, user := range users.Users {
		if strings.Compare(user_name, user.BrodcasterLogin) == 0 {
			u = user
		}
	}

	return &u, nil
}
