package controllers

import (
	"database/sql"
	"github.com/hugorjs/url-shortner-app/internal/db"
	"github.com/hugorjs/url-shortner-app/internal/url"
	"html/template"
	"net/http"
	"strings"
)

func Shorten(lite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if err := db.StoreURL(lite, shortUrl, originalURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
}

func Proxy(lite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := r.URL.Path[1:]
		if shortUrl == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
		origUrl, err := db.GetOriginalURL(lite, shortUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, origUrl, http.StatusPermanentRedirect)
	}
}
