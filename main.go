package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	Id           string    `json:"id"`
	OriginalURL  string    `json:"originalURL"`
	ShortenURL   string    `json:"shortenURL"`
	CreationDate time.Time `json:"creationTime"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:6] // Only first 6 characters
}

func createUrl(OriginalURL string) string {
	shortUrl := generateShortURL(OriginalURL)
	id := shortUrl
	urlDB[id] = URL{
		Id:           id,
		OriginalURL:  OriginalURL,
		ShortenURL:   shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

func getUrl(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, fmt.Errorf("URL not found")
	}
	return url, nil
}

func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"` // Correct field name here
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := createUrl(data.URL)
	response := struct {
		ShortURL string `json:"shortURL"`
	}{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):] // Extract short URL ID
	url, err := getUrl(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
