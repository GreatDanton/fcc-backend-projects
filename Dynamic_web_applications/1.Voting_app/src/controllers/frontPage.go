package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
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
	LastID       int
	Pagination   bool
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
	case "POST":
		nextPage(w, r)
	default:
		displayFrontPage(w, r)
	}
}

// getMaxID from url that defines pool with maximum id parsed from db
func getMaxID(r *http.Request) (int, error) {
	q := r.URL.Query()
	urlID := q.Get("maxID")
	maxID := 0
	if urlID != "" {
		id, err := strconv.Atoi(urlID)
		if err != nil {
			return maxID, err
		}
		maxID = id
	}
	return maxID, nil
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
	maxID, err := getMaxID(r)
	// user added something into url
	if err != nil {
		fmt.Println(err)
	}

	// getting database response based on the maxID
	rows, err := fpQuery(maxID)

	if err != nil {
		log.Fatal(err)
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

	// grab last id
	lastID, err := strconv.Atoi(pools[len(pools)-1].ID)
	if err != nil {
		fmt.Println(err)
	}

	user := LoggedIn(r)

	// display next button or not?
	dp := true
	if len(pools) < 20 { // if array of posts is less than limit of our query
		dp = false
	}

	fp := frontPage{Pools: pools, LoggedInUser: user, LastID: lastID, Pagination: dp}

	// displaying template
	t := template.Must(template.ParseFiles("templates/frontPage.html", "templates/navbar.html", "templates/pagination.html"))
	err = t.ExecuteTemplate(w, "frontPage", fp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

// nextPage paginates the results and displays the next x pools
func nextPage(w http.ResponseWriter, r *http.Request) {

}
