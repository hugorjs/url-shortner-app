package controllers

import (
	"github.com/hugorjs/url-shortner-app/internal/url"
	"html/template"
	"net/http"
	"strings"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL not provided", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "http://" + originalURL
	}

	// shorten then URL.
	shortUrl := url.Shorten(originalURL)

	data := map[string]string{
		"ShortURL": shortUrl,
	}

	t, err := template.ParseFiles("internal/views/shorten.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
