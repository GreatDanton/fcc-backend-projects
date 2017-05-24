package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	fmt.Println("Starting server on", "http://127.0.0.1:8080")

	http.HandleFunc("/api/imagesearch/", search)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", "Search part")
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
	fmt.Println(server)

	url := fmt.Sprintf("%v", server)
	str, _ := createImageAPI(url)
	fmt.Fprintf(w, "%s\n", str)
	/*	fmt.Println(str)*/
}

// struct for outputting json
type output struct {
	URL       string `json:"url"`
	Snippet   string `json:"snippet"`
	Thumbnail string `json:"thumbnail"`
	Context   string `json:"context"`
}

// parses html
func createImageAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("Could not parse response html: %v", err)
		return "", e
	}

	return string(text), nil
}
