package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", rootHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	fmt.Println("hello world")
}

type jsonOutput struct {
	Unix    string `json:"unix"`
	Natural string `json:"natural"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Path[1:]

	data := formatDate(input)
	out, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Json marshal failed: %s", err)
	}
	fmt.Fprintf(w, "%s\n", out)
}

// formatting date out of input string
func formatDate(input string) jsonOutput {
	data := jsonOutput{Unix: "", Natural: ""}
	// check for unix format
	if len(input) > 0 {
		inputArr := strings.Split(input, ",")

		switch len(inputArr) {
		case 1:
			u, err := strconv.ParseInt(inputArr[0], 10, 64)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return data
			}
			t := time.Unix(u, 0)
			fmt.Println(t)
			data.Unix = inputArr[0]
			data.Natural = t.Format("January 2, 2006")
		case 2:
			fmt.Println("bla bla")

		default:
			return data
		}
	}

	return data
}
