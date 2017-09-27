package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

type frontPage struct {
	Pools        []pool
	LoggedInUser User
	Pagination   pagination
	/* 	MaxIDNext      string
	   	MaxIDPrev      string
	   	PaginationNext bool
	   	PaginationPrev bool */
}

type pagination struct {
	MaxIDNext      string
	MaxIDPrev      string
	PaginationNext bool
	PaginationPrev bool
}

type pool struct {
	ID         string
	Time       string
	Title      string
	Author     string
	NumOfVotes string
}

// FrontPage takes care of displaying front page of Voting Application
func FrontPage(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		displayFrontPage(w, r)
	default:
		displayFrontPage(w, r)
	}
}

// getMaxID from url that defines pool with maximum id parsed from db
func getURLParams(r *http.Request) int {
	q := r.URL.Query()
	urlID := q.Get("maxID")
	maxID := 0
	if urlID != "" {
		id, err := strconv.Atoi(urlID)
		// user added something into url
		if err != nil {
			return maxID
		}
		maxID = id
	}
	return maxID
}

// fpQuery picks the most suitable sql query based on maxID of pool and returns
// sql rows and error
func fpQuery(maxID int) (*sql.Rows, error) {
	if maxID == 0 { // maxID = 0 when we perform the first query on "/"
		rows, err := global.DB.Query(`SELECT pool.id, pool.title,
								  users.username, pool.time,
								  (select count(*) as votes from vote where vote.pool_id = pool.id)
								  FROM pool
								  LEFT JOIN users on users.id = pool.created_by
								  ORDER BY pool.id desc
								  limit 20`)
		return rows, err
	}

	rows, err := global.DB.Query(`SELECT pool.id, pool.title,
								  users.username, pool.time,
								  (select count(*) as votes from vote where vote.pool_id = pool.id)
								  FROM pool
								  LEFT JOIN users on users.id = pool.created_by
								  WHERE pool.id <= $1
								  ORDER BY pool.id desc
								  limit 20`, maxID)
	return rows, err
}

// displaysFrontPage with latest pools
func displayFrontPage(w http.ResponseWriter, r *http.Request) {
	maxID := getURLParams(r)

	// getting database response based on the maxID
	rows, err := fpQuery(maxID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var (
		id         string
		title      string
		author     string
		time       time.Time
		numOfVotes string
	)

	pools := []pool{}

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &time, &numOfVotes)
		if err != nil {
			log.Fatal(err)
		}
		// get time difference in human readable format
		t := utilities.TimeDiff(time)
		pools = append(pools, pool{ID: id, Title: title,
			Author: author, Time: t, NumOfVotes: numOfVotes})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	user := LoggedIn(r)
	p := handlePoolPagination(maxID, pools)
	fp := frontPage{Pools: pools, LoggedInUser: user, Pagination: p}

	// displaying template
	err = global.Templates.ExecuteTemplate(w, "frontPage", fp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

// displayNextButton returns newMaxID for displaying it in url and bool (true/false)
// true => display next button
// false => do not display next button
func displayNextButton(pools []pool) (string, bool) {
	newMaxID := pools[len(pools)-1].ID
	displayPagination := true
	if len(pools) < 20 {
		displayPagination = false
	}
	return newMaxID, displayPagination
}

// displayPrevButton determine wheter to display previous button or not
// and returns MaxIDPrev for displaying previous page and bool
// true => display previous button
// false => hide previous button
func displayPrevButton(maxID int, pools []pool) (string, bool) {
	// when maxID = 0, we are on the front page
	if maxID == 0 {
		return "", false
	}
	// if maxID is bigger than first item in pools that means we are on the front
	// page and should not display previous button
	poolID, err := strconv.Atoi(pools[0].ID)
	if err != nil {
		return "", false
	}
	// this means we are coming from previous
	// to the front page via prev button
	if maxID > poolID { //If previous maxID + 20 > current pool.ID
		return "", false
	}
	// if none of the above applies move to previous page by increasing
	// maxID by 20
	maxID += 20
	id := strconv.Itoa(maxID)
	return id, true
}

// returns pagination struct that handles moving back and forth between
// the pool pages and displaying appropriate buttons
func handlePoolPagination(maxID int, pools []pool) pagination {
	maxIDNext, dpn := displayNextButton(pools)
	maxIDPrev, dpp := displayPrevButton(maxID, pools)
	p := pagination{MaxIDNext: maxIDNext, PaginationNext: dpn,
		MaxIDPrev: maxIDPrev, PaginationPrev: dpp}
	return p
}
