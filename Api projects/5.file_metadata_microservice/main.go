package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// template for front page as a global variable, parsing template only at the startup
var frontPageTemplate = template.Must(template.ParseFiles("templates/index.html"))

// main function starting server
func main() {
	fmt.Println("Starting server on", "http://127.0.0.1:8080")

	http.HandleFunc("/", indexPage)
	// use public folder
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

// fileAPI is used for displaying json string of file metadata
type fileAPI struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// handler for front page of the microservice
func indexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// display tempalte on get request
		t := frontPageTemplate
		err := t.Execute(w, "")
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		// parse file and header from upload-file input field
		file, header, err := r.FormFile("upload-file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// get filesize from file
		var buff bytes.Buffer
		fileSize, err := buff.ReadFrom(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		api := fileAPI{Name: header.Filename, Size: fileSize}

		// create json byte string
		json, err := json.MarshalIndent(api, "", "    ")
		if err != nil {
			fmt.Println(err)
			return
		}
		// display json string in human readable format
		fmt.Fprintf(w, fmt.Sprintf("%s", json))
	}
}
