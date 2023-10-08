package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"twitch_web_app/models"
)

func view(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	q := r.URL.Query()
	search := q.Get("search")
	names := strings.Fields(search)
	var chans []<-chan models.User

	for _, name := range names {
		chans = append(chans, find_streamer(name))
	}
	for n := range merge(chans...) {
		users = append(users, n)
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

func find_streamer(name string) <-chan models.User {
	out := make(chan models.User, 1)
	go func() {
		defer close(out)
		data, _ := models.Search(name)
		out <- *data
	}()
	return out
}

func merge(cs ...<-chan models.User) <-chan models.User {
	var wg sync.WaitGroup
	out := make(chan models.User)

	output := func(c <-chan models.User) {
		defer wg.Done()

		for n := range c {
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
