package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (o *Rest) HShortenerURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		o.HShortener(w, r)
	} else if r.Method == http.MethodGet {
		o.HRedirect(w, r)
	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func (o *Rest) HShortener(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	reqObj := strings.TrimSpace(string(body))

	shortURL := o.generateShortURL()

	o.urlMap[shortURL] = reqObj

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortURL)
}

func (o *Rest) HRedirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.RequestURI[1:]

	URL := o.urlMap[shortURL]

	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, r, URL, http.StatusTemporaryRedirect)
}
