package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// pagination handles pagination display
type pagination struct {
	MaxIDNext      string
	MaxIDPrev      string
	URLPath        string
	PaginationNext bool
	PaginationPrev bool
}

// handlePollPagination returns pagination struct that handles moving back and forth between
// the poll pages and displaying appropriate buttons
func handlePollPagination(r *http.Request, maxID int, polls []poll, limit int) pagination {
	urlPath := r.URL.EscapedPath()
	urlPath = strings.TrimRight(urlPath, "/") // remove right "/"
	maxIDNext, dpn := displayNextButton(polls, limit)
	maxIDPrev, dpp := displayPrevButton(maxID, polls, limit)

	p := pagination{MaxIDNext: maxIDNext, PaginationNext: dpn, MaxIDPrev: maxIDPrev, PaginationPrev: dpp, URLPath: urlPath}
	return p
}

// displayNextButton returns newMaxID for displaying it in url and bool (true/false)
// true => display next button
// false => do not display next button
func displayNextButton(polls []poll, limit int) (string, bool) {
	displayPagination := true
	if len(polls) < limit {
		displayPagination = false
		return "", displayPagination
	}
	newMaxID := polls[len(polls)-1].ID // last item in polls id array
	return newMaxID, true
}

// displayPrevButton determine wheter to display previous button or not
// and returns MaxIDPrev for displaying previous page and bool
// true => display previous button
// false => hide previous button
func displayPrevButton(maxID int, polls []poll, limit int) (string, bool) {
	// if polls do not exist, do not display previous button
	if len(polls) == 0 {
		return "", false
	}
	// when maxID = 0, we are on the front page
	if maxID == 0 {
		return "", false
	}
	// if maxID is bigger than first item in polls that means we are on the front
	// page and should not display previous button
	currentMaxPollID, err := strconv.Atoi(polls[0].ID)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	// coming from previous to front page
	if maxID > currentMaxPollID { // this means we are already on the front page
		return "", false
	}
	// if none of the above applies move to previous page by increasing
	// maxID by limit
	maxID += limit
	id := strconv.Itoa(maxID)
	return id, true
}
