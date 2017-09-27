package controllers

import (
	"database/sql"
	"fmt"
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
}

// Pool is structure that contains all relevant pool data
type pool struct {
	ID         string
	Time       string
	Title      string
	Author     string
	NumOfVotes string
}

// FrontPage takes care of displaying front page of Voting Application
func FrontPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayFrontPage(w, r)
	default:
		displayFrontPage(w, r)
	}
}

// getMaxID from url that defines pool with maximum id parsed from db
func getMaxIDParam(r *http.Request) int {
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

// displaysFrontPage with latest pools
func displayFrontPage(w http.ResponseWriter, r *http.Request) {
	maxID := getMaxIDParam(r)
	limit := 10
	// getting database response based on the maxID
	pools, err := getFrontPageData(maxID, limit)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := LoggedIn(r)
	p := handlePoolPagination(maxID, pools, limit)
	fp := frontPage{Pools: pools, LoggedInUser: user, Pagination: p}

	// displaying template
	err = global.Templates.ExecuteTemplate(w, "frontPage", fp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetFrontPageData returns array of pools based on chosen maxID(max pool id) and limit of results
func getFrontPageData(maxID int, limit int) ([]pool, error) {
	pools := []pool{}
	rows, err := fpQuery(maxID, limit)
	if err != nil {
		return pools, err
	}
	defer rows.Close()

	var (
		id         string
		title      string
		author     string
		time       time.Time
		numOfVotes string
	)

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &time, &numOfVotes)
		if err != nil {
			return pools, err
		}
		// get time difference in human readable format
		t := utilities.TimeDiff(time)
		pools = append(pools, pool{ID: id, Title: title,
			Author: author, Time: t, NumOfVotes: numOfVotes})
	}
	err = rows.Err()
	if err != nil {
		return pools, err
	}
	// if error does not happen, return results
	return pools, nil
}

// fpQuery picks the most suitable sql query based on maxID of pool and returns
// sql rows and error
func fpQuery(maxID int, limit int) (*sql.Rows, error) {
	if maxID == 0 { // maxID = 0 when we perform the first query on "/"
		rows, err := global.DB.Query(`SELECT pool.id, pool.title,
								  users.username, pool.time,
								  (select count(*) as votes from vote where vote.pool_id = pool.id)
								  FROM pool
								  LEFT JOIN users on users.id = pool.created_by
								  ORDER BY pool.id desc
								  limit $1`, limit)
		return rows, err
	}

	rows, err := global.DB.Query(`SELECT pool.id, pool.title,
								  users.username, pool.time,
								  (select count(*) as votes from vote where vote.pool_id = pool.id)
								  FROM pool
								  LEFT JOIN users on users.id = pool.created_by
								  WHERE pool.id < $1
								  ORDER BY pool.id desc
								  limit $2`, maxID, limit)
	return rows, err
}
