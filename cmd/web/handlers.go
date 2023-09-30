package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"twitch_web_app/models"
)

func view(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	q := r.URL.Query()
	search := q.Get("search")
	names := strings.Fields(search)

	for _, name := range names {
		user_info, err := models.Search(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		users = append(users, *user_info)
	}

	ts, err := template.ParseFiles("./ui/html/pages/view.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ts.Execute(w, users)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/pages/home.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
