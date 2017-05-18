package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/", rootHandle)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listen and serve:", err)
	}
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	// parse language and user-agent
	header := r.Header
	addr := r.RemoteAddr

	// get os
	h := fmt.Sprintf("%v", header)
	os := parseOs(h)
	fmt.Println(os)

	// get language
	l := parseLang(h)
	fmt.Println(l)

	// get ip
	ip := parseIP(addr)
	fmt.Println(ip)

	fmt.Fprintf(w, "%s\n", header)
	fmt.Fprintf(w, "%v\n", addr)
}

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
