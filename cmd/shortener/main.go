package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	base62 "github.com/deatil/go-encoding/base62"
)

var db = make(map[string]string, 0)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", ShortenUrl)
	// mux.HandleFunc("/{id}", GetUrl)
	http.HandleFunc("/", Handler)
	return http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		body, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		encodedUrl := base62.StdEncoding.EncodeToString(body)
		// Записать в бд
		db[encodedUrl] = string(body)
		w.WriteHeader(201)
		w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", encodedUrl)))

	} else if req.Method == http.MethodGet {
		
		id := strings.TrimPrefix(req.URL.Path, "/")
		v, ok := db[id]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", v)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}