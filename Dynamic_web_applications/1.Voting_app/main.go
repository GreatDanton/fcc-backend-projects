package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server: http://127.0.0.1:8080")

	http.HandleFunc("/", frontPage)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
