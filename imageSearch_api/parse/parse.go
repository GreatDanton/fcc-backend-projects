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
	Snippet   string `json:"snippet"`
	Thumbnail string `json:"thumbnail"`
	Context   string `json:"context"`
}

// CreateImageAPI from provided url string
func CreateImageAPI(url string) ([]byte, error) {
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

	/*	fmt.Println("### Printing meta tag")
		fmt.Println(meta)*/
	return meta, nil
}

func parseHTML(html string) (string, error) {
	return html, nil
}

// TODO: remove this part
// used for parsing json from string on google images
/*type parseMetadata struct {
	Description string `json:"pt"`
	Image       string `json:"ou"`
	Site        string `json:"ru"`
}*/

// parses relevant metadata from html string
// returns: json api
func getMetadata(html string) ([]byte, error) {

	api := []ImageAPI{}

	imageContainers := parseImageContainers(html, []string{})

	for _, container := range imageContainers {
		image := imageURL(container)
		context := imageContext(container)
		img := ImageAPI{Context: context, Thumbnail: image}
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
	divEnd := strings.Index(html[divStart+1:], `</td>`) + divStart

	if divStart == -1 {
		return containers
	}

	imgContainer := html[divStart : divEnd+1]
	containers = append(containers, imgContainer)

	return parseImageContainers(html[divEnd:], containers)
}

// get image url from provided html string
func imageURL(html string) string {
	imgStart := strings.Index(html, `<img`)
	h := html[imgStart:]

	start := strings.Index(h, `src`)
	h = h[start+len(`src="`):]
	end := strings.Index(h, `"`)

	return h[:end]
}

// get image context from provided html string
func imageContext(html string) string {
	cStart := strings.Index(html, `</cite>`)
	h := html[cStart+len(`</cite><br>`):]
	end := strings.Index(h, `<br>`)

	context := h[:end]
	// crude way to replace <b> </b> tags from context string
	context = strings.Replace(context, `<b>`, "", -1)
	context = strings.Replace(context, `</b>`, "", -1)

	return context
}
