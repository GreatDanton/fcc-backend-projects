package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", mainHandle)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

// reply struct for marshalling output data into json
type reply struct {
	Language string `json:"language"`
	Os       string `json:"software"`
	IP       string `json:"ipaddress"`
}

// mainHandle function
func mainHandle(w http.ResponseWriter, r *http.Request) {
	// parse language and user-agent
	header := r.Header
	h := fmt.Sprintf("%v", header)
	addr := r.RemoteAddr

	// get os
	os := parseOs(h)
	// get language
	l := parseLang(h)
	// get ip
	ip := parseIP(addr)

	reply := reply{Language: l, Os: os, IP: ip}

	out, err := json.MarshalIndent(reply, "", "    ")
	if err != nil {
		fmt.Println("Cannot produce json", err)
		return
	}

	// output json
	fmt.Fprintf(w, "%s\n", out)
}

//
// HELPER FUNCTIONS
//

// parse ip from remote address
func parseIP(addr string) string {
	ip := strings.Split(addr, ":")[0]
	return ip
}

// parse language from header
func parseLang(h string) string {
	tag := strings.Index(h, "Accept-Language:")
	start := strings.Index(h[tag:], "[")
	start += tag + 1

	end := strings.Index(h[start:], "]")
	end += start

	l := h[start:end]

	return strings.Split(l, ";")[0]
}

// parse of from header string
func parseOs(h string) string {
	tag := strings.Index(h, "User-Agent")
	start := strings.Index(h[tag:], "(")
	start += tag + 1 // +1 to remove the first bracket

	end := strings.Index(h[start:], ")")
	end += start

	return h[start:end]
}
