package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
)

type URLCreationRequest struct {
	LongURL string `json:"long_url" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

type URLCreationResponse struct {
	Message  string `json:"message"`
	ShortURL string `json:"short_url"`
}

type saveData struct {
	userId   string
	longURL  string
	shortURL string
}

var ctx = context.Background()

func main() {
	rdb := store.Redis()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			d := json.NewDecoder(r.Body)

			var data URLCreationRequest
			if err := d.Decode(&data); err != nil {
				panic(err)
			}

			shortURL := GenerateShortLink(data.LongURL, data.UserID)

			payload := URLCreationResponse{
				Message:  "short url created successfully",
				ShortURL: "http://localhost:3000/" + shortURL,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(payload)

			err := rdb.Set(ctx, shortURL, data.LongURL, 0).Err()
			if err != nil {
				panic(err)
			}

		} else if r.Method == "GET" {
			s := strings.TrimPrefix(r.URL.Path, "/")

			redirectURL, err := rdb.Get(ctx, s).Result()
			if err == redis.Nil {
				fmt.Fprintf(w, "URL not found")
			} else if err != nil {
				panic(err)
			}

			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
