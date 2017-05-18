package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
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

		in := strings.Replace(input, ",", "", -1)
		inputArr := strings.Split(in, " ")

		// only unix time is posted in url
		if len(inputArr) == 1 {
			u, err := strconv.ParseInt(inputArr[0], 10, 64)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				return data
			}
			t := time.Unix(u, 0)
			data.Unix = inputArr[0]
			data.Natural = t.Format("January 2, 2006")
			return data
		}

		// if there are more words => full date is posted
		// we are parsing in american format: MM DD YYYY
		y := ""
		d := ""
		m := ""

		mPlace := monthNameCheck(inputArr)
		if mPlace >= 0 { // month name exists in query
			m = inputArr[mPlace]
			restArr := inputArr
			// removing element from slice
			restArr = append(restArr[:mPlace], restArr[mPlace+1:]...) // lol wut?
			yPlace := maxLen(restArr)
			switch yPlace {
			case 0:
				y = restArr[0]
				d = restArr[1]
			case 1:
				y = restArr[1]
				d = restArr[0]
			}

		} else {
			// month name does not exist
			yPlace := maxLen(inputArr)
			switch yPlace {
			case 0: // yyyy mm dd
				y = inputArr[0]
				m = inputArr[1]
				d = inputArr[2]
			case 1: // mm yyyy dd
				y = inputArr[1]
				m = inputArr[0]
				d = inputArr[2]
			case 2: // mm dd yyyy
				y = inputArr[2]
				m = inputArr[0]
				d = inputArr[1]
			}
		}

// transforming data to integers
		year, err := strconv.Atoi(y)
		if err != nil {
			fmt.Println("Problem with converting year:", err)
			return data
		}

		month := 0
		// month
		// if month (m) starts with letter, get the integer
		if isLetter(m) {
			mArr := []string{"---", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
			for i, v := range mArr {
				if strings.Contains(v, m) {
					month = i
					break
				}
			}
		} else { // if month (m) is an integer
			month, err = strconv.Atoi(m)
			if err != nil {
				fmt.Println("Problem with converting month:", err)
				return data
			}
		}

		day, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Problem with converting day:", err)
			return data
		}

		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		data.Unix = fmt.Sprintf("%v", t.Unix())
		data.Natural = fmt.Sprintf("%v", t.Format("January 2, 2006"))
	}

	return data
}

func maxLen(arr []string) int {
	index := 0
	maxLen := 0
	for i, v := range arr {
		if len(v) > maxLen {
			maxLen = len(v)
			index = i
		}
	}
	return index
}

// check if string contains letters
// usage: isLetter("Something") => true
// isLetter("123") => false
var isLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

//
// check if string / month name exists in array
// returns index
// or -1 if it doesn't exist
func monthNameCheck(arr []string) int {
	for i, v := range arr {
		if isLetter(v) {
			return i
		}
	}
	return -1
}
