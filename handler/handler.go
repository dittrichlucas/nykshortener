package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/dittrichlucas/nykshortener/shortener"
	"github.com/dittrichlucas/nykshortener/store"
	"github.com/google/uuid"
)

type URLCreationRequest struct {
	LongURL string `json:"long_url" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

type URLCreationResponse struct {
	Message  string `json:"message"`
	ShortURL string `json:"short_url"`
}

func CreateShortURL(w http.ResponseWriter, r *http.Request) URLCreationResponse {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	longURL := r.PostForm.Get("shortener")

	userId := uuid.New().Version().String()
	shortURL := shortener.GenerateShortLink(longURL, userId)

	log.Println("--> userId: " + userId)
	log.Println("--> shortURL: " + shortURL)

	payload := URLCreationResponse{
		Message:  "short url created successfully",
		ShortURL: "http://localhost:3000/" + shortURL,
	}

	store.SaveURLMapping(shortURL, longURL)

	return payload
}

func ShortURLRedirect(w http.ResponseWriter, r *http.Request) {
	s := strings.TrimPrefix(r.URL.Path, "/")

	redirectURL := store.RetrieveInitialURL(s)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
