package parse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ImageAPI for displaying image json api
type ImageAPI struct {
	URL       string `json:"url"`
	Snippet   string `json:"snippet"`
	Thumbnail string `json:"thumbnail"`
	Context   string `json:"context"`
}

// CreateImageAPI from provided url string
func CreateImageAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("Could not parse response html: %v", err)
		return "", e
	}
	html := string(body)
	// get url, snippet, thumbnail, context
	meta, _ := getMetadata(html)
	fmt.Println(meta)

	return html, nil
}

func parseHTML(html string) (string, error) {
	return html, nil
}

// used for parsing json from string on google images
type parseMetadata struct {
	Description string `json:"pt"`
	Image       string `json:"ou"`
	Site        string `json:"ru"`
}

func getMetadata(html string) (string, error) {
	divStart := strings.Index(html, `<td style="width:25%`)
	divEnd := strings.Index(html[divStart+1:], `</td>`) + divStart
	fmt.Println(divStart)

	if divStart == -1 {
		err := fmt.Errorf("No more meta divs")
		return "", err
	}

	image := html[divStart : divEnd+1]

	// json string we will unmarshal
	/*	j := []byte(html[divStart:divEnd])

		api := parseMetadata{}
		err := json.Unmarshal(j, api)
		if err != nil {
			return parseMetadata{}, fmt.Errorf("getMetadata: error parsing json %v", err)
		}

		return api, nil*/
	return image, nil
}
