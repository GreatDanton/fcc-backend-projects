package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/imageSearch_api/parse"
)

// SearchStorage global search storage
var SearchStorage []string

func main() {
	fmt.Println("Starting server on", "http://127.0.0.1:8080")

	http.HandleFunc("/api/imagesearch/", search)
	http.HandleFunc("/", frontPageHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	serverURL := "https://www.google.com/search?tbm=isch"

	// parse queries from url
	offset, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return
	}

	// get search query from url
	searchQuery := strings.Split(r.URL.Path, "/")[3:][0]
	fmt.Println(searchQuery)

	server, err := url.Parse(serverURL)
	if err != nil {
		log.Fatal("/imagesearch/ error parsing serverURL, error:", err)
	}

	// set searching query to our search string
	q := server.Query()
	q.Set("q", searchQuery)

	// add offset to google image search -> original offset multiplied by 10
	if val, ok := offset["offset"]; ok {
		off, err := strconv.Atoi(val[0])
		if err != nil {
			fmt.Println("Cannot turn offset string into integer")
			return
		}
		q.Set("start", fmt.Sprint(off*10))
	}

	server.RawQuery = q.Encode()
	url := fmt.Sprintf("%v", server)

	json, err := parse.CreateImageAPI(url)
	if err != nil {
		log.Fatal("Cannot create image api", err)
	}
	fmt.Fprintf(w, "%s\n", json)

	// on each search add searched url to global SearchStorage
	if len(SearchStorage) < 10 {
		SearchStorage = append(SearchStorage, searchQuery)
	} else {
		SearchStorage = SearchStorage[1:] // remove first element from array
		SearchStorage = append(SearchStorage, searchQuery)
	}
}

// handling main page of the api service
func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "use /api/imagesearch/[search query] to search for image\n\n")
	fmt.Fprintf(w, "Last searched for: \n")
	for _, i := range SearchStorage {
		fmt.Fprintf(w, fmt.Sprintf("	%v \n", i))
	}
}
