package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/imageSearch_api/parse"
)

func main() {
	fmt.Println("Starting server on", "http://127.0.0.1:8080")

	http.HandleFunc("/api/imagesearch/", search)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	serverURL := "https://www.google.com/search?tbm=isch"
	// remove /api/imagesearch/ from query
	search := strings.Split(r.URL.Path, "/")[3:][0]
	fmt.Println(search)

	server, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("/imagesearch/ error parsing serverURL, error:", err)
	}

	// set searching query to our search string
	q := server.Query()
	q.Set("q", search)
	server.RawQuery = q.Encode()
	url := fmt.Sprintf("%v", server)

	json, err := parse.CreateImageAPI(url)
	if err != nil {
		log.Fatal("Cannot create image api", err)
	}
	fmt.Fprintf(w, "%s\n", json)
}
