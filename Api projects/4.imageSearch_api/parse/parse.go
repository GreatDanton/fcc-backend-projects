package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ImageAPI used for displaying image json api in browser
type ImageAPI struct {
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Context   string `json:"context"`
	// Snippet   string `json:"snippet"`
}

// CreateImageAPI from provided url string
func CreateImageAPI(url string) ([]byte, error) {
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("Could not parse response html: %v", err)
		return nil, e
	}
	html := string(body)
	// get url, snippet, thumbnail, context
	meta, err := getMetadata(html)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// parses relevant metadata from html string
// returns: json api
func getMetadata(html string) ([]byte, error) {

	api := []ImageAPI{}

	imageContainers := parseImageContainers(html, []string{})

	for _, container := range imageContainers {
		image := parseThumbnailURL(container)
		context := parseImageContext(container)
		url := parseSiteURL(container)
		img := ImageAPI{Context: context, Thumbnail: image, URL: url}
		api = append(api, img)
	}

	jsonString, err := json.MarshalIndent(api, "", "    ")
	if err != nil {
		fmt.Println("getMetadata: Cannot marshal image api")
		return nil, err
	}

	return jsonString, nil
}

// recursive function for parsing image parent containers from raw
// html string
func parseImageContainers(html string, containers []string) []string {
	divStart := strings.Index(html, `<td style="width:25%`)
	divEnd := strings.Index(html[divStart+1:], `</td>`) + divStart + len(`</td>`)

	if divStart == -1 {
		return containers
	}

	imgContainer := html[divStart : divEnd+1]
	/* 	fmt.Println("")
	   	fmt.Println(imgContainer) */
	containers = append(containers, imgContainer)
	return parseImageContainers(html[divEnd:], containers)
}

// returns website on which the image is found - url part of json api
func parseSiteURL(html string) string {
	imgStart := strings.Index(html, `<a href`)
	html = html[imgStart:]
	s := strings.Index(html, "http")
	end := strings.Index(html[s:], `&amp`) + s
	url := html[s:end]
	return url
}

// get image url from provided html string
func parseThumbnailURL(html string) string {
	imgStart := strings.Index(html, `<img`)
	h := html[imgStart:]

	start := strings.Index(h, `src`)
	h = h[start+len(`src="`):]
	end := strings.Index(h, `"`)

	return h[:end]
}

// get image context from provided html string
func parseImageContext(html string) string {
	cStart := strings.Index(html, `</cite>`)
	h := html[cStart+len(`</cite><br>`):]
	end := strings.Index(h, `<br>`)

	context := h[:end]
	// crude way to replace <b> </b> tags from context string
	context = strings.Replace(context, `<b>`, "", -1)
	context = strings.Replace(context, `</b>`, "", -1)

	return context
}
