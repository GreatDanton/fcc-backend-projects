package controllers

import (
	"strconv"
)

// pagination handles pagination display
type pagination struct {
	MaxIDNext      string
	MaxIDPrev      string
	PaginationNext bool
	PaginationPrev bool
}

// handlePoolPagination returns pagination struct that handles moving back and forth between
// the pool pages and displaying appropriate buttons
func handlePoolPagination(maxID int, pools []pool, count int) pagination {
	maxIDNext, dpn := displayNextButton(pools, count)
	maxIDPrev, dpp := displayPrevButton(maxID, pools, count)
	p := pagination{MaxIDNext: maxIDNext, PaginationNext: dpn,
		MaxIDPrev: maxIDPrev, PaginationPrev: dpp}
	return p
}

// displayNextButton returns newMaxID for displaying it in url and bool (true/false)
// true => display next button
// false => do not display next button
func displayNextButton(pools []pool, count int) (string, bool) {
	newMaxID := pools[len(pools)-1].ID
	displayPagination := true
	if len(pools) < count {
		displayPagination = false
	}
	return newMaxID, displayPagination
}

// displayPrevButton determine wheter to display previous button or not
// and returns MaxIDPrev for displaying previous page and bool
// true => display previous button
// false => hide previous button
func displayPrevButton(maxID int, pools []pool, count int) (string, bool) {
	// when maxID = 0, we are on the front page
	if maxID == 0 {
		return "", false
	}
	// if maxID is bigger than first item in pools that means we are on the front
	// page and should not display previous button
	maxPoolID, err := strconv.Atoi(pools[0].ID)
	if err != nil {
		return "", false
	}
	// this means we are coming from previous to the front page via prev button
	//If previous maxID + count > current pool.ID
	// adding +1 to maxPoolID, othwerwise the if would be always true
	if maxID > maxPoolID+1 {
		return "", false
	}
	// if none of the above applies move to previous page by increasing
	// maxID by 20
	maxID += count
	id := strconv.Itoa(maxID)
	return id, true
}
