package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	godotenv.Load(".env")
	get_auth()

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", home)
	router.HandlerFunc(http.MethodGet, "/view", view)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func get_auth() error {
	url := "https://id.twitch.tv/oauth2/token?client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET") + "&grant_type=client_credentials"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var b map[string]string
	json.Unmarshal(body, &b)
	os.Setenv("BEARER_TOKEN", b["access_token"])
	return nil

}
