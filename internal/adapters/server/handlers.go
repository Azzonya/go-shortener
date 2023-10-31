package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (o *Rest) HShortener(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	reqObj := string(body)

	shortUrl := o.generateShortURL()

	o.urlMap[shortUrl] = reqObj

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, shortUrl)
}

func (o *Rest) HRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	shortUrl := queryParams.Get("id")

	URL := o.urlMap[shortUrl]

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Fprint(w, URL)
}
