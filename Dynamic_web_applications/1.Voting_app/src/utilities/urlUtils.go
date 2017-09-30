package utilities

import (
	"net/http"
	"strings"
)

// GetURLSuffix returns last part of the url
// url: poll/10/edit => returned: edit
func GetURLSuffix(r *http.Request) string {
	url := strings.TrimRight(r.URL.Path, "/")
	path := strings.Split(url, "/")
	suffix := path[len(path)-1]
	return suffix
}
