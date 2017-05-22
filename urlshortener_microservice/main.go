package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// database storing [number] = url
var db = make(map[string]string)
var index = template.Must(template.ParseFiles("index.html"))

type indexTemp struct {
	HOST string
}

func main() {
	http.HandleFunc("/new/", newLink)
	http.HandleFunc("/", mainHandler)

	fmt.Println("Server started on: http://127.0.0.1:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

// func for handling main website
func mainHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Path
	// if url is root url, display root template
	if u == "/" {
		t := index
		err := t.Execute(w, indexTemp{HOST: r.Host})
		if err != nil {
			fmt.Println("Error executing template")
		}
		return
	}

	url := u[1:]
	site, ok := db[url]
	if !ok {
		fmt.Fprintf(w, "%s\n", "Link does not exist")
		displayDB(w)
		return
	}
	// if redirection is not working on local host it might
	// be because of the browser cache -> test it in private mode
	http.Redirect(w, r, site, http.StatusSeeOther)
}

// struct for json output on new link creation
type newLinkOutput struct {
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}

// for displaying error message in json
type errOut struct {
	Error string `json:"error"`
}

// function for handling new links
func newLink(w http.ResponseWriter, r *http.Request) {
	u := r.URL.RequestURI()[len("/new/"):] // remove /new/ part from url
	// add / between http and url
	// parsed from url => http:/something
	// valid url       => http://something
	httpPart := strings.Index(u, "/")
	url := u[:httpPart] + "//" + u[httpPart+1:]
	host := fmt.Sprintf("%s/", r.Host)
	newLink := ""

	// if http is not found in url return
	if strings.Index(url, "http") == -1 {
		e := errOut{Error: "Wrong url format, make sure you have a valid protocol and link"}
		out, err := json.Marshal(e)
		if err != nil {
			fmt.Println("Error on json marshal:", err)
			return
		}
		fmt.Fprintf(w, "%s\n", out)
		displayDB(w)
		return
	}

	for {
		s1 := rand.NewSource(time.Now().UnixNano()) // new seed
		randFloat := rand.New(s1).Float64()         // random number

		newLink = fmt.Sprintf("%.5f", randFloat)[2:]
		_, ok := db[newLink] // ok == false if it does not exist
		if !ok {             // if false (does not exist) add url to db
			db[newLink] = url
			break
		}
	}

	outStruct := newLinkOutput{Original: url, Short: host + newLink}
	out, err := json.MarshalIndent(outStruct, "", "    ")
	if err != nil {
		fmt.Println("NewLink, error marshalling json", err)
		return
	}

	fmt.Fprintf(w, "%s\n", out)
	fmt.Fprint(w, "\n")
}

// display whole database -> obviously not suitable for production
func displayDB(w http.ResponseWriter) {
	jsonString, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		fmt.Println("displayDB err json marshalling;", err)
		return
	}
	fmt.Fprintf(w, "%s\n", jsonString)
}
