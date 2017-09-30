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
	Polls        []poll
	LoggedInUser User
	Pagination   pagination
}

// Poll is structure that contains all relevant poll data
type poll struct {
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

// getMaxID from url that defines poll with maximum id parsed from db
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

// displaysFrontPage with latest polls
func displayFrontPage(w http.ResponseWriter, r *http.Request) {
	maxID := getMaxIDParam(r)
	limit := 20
	// getting database response based on the maxID
	polls, err := getFrontPageData(maxID, limit)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := LoggedIn(r)
	p := handlePollPagination(r, maxID, polls, limit)
	fp := frontPage{Polls: polls, LoggedInUser: user, Pagination: p}

	// displaying template
	err = global.Templates.ExecuteTemplate(w, "frontPage", fp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetFrontPageData returns array of polls based on chosen maxID(max poll id) and limit of results
func getFrontPageData(maxID int, limit int) ([]poll, error) {
	polls := []poll{}
	rows, err := fpQuery(maxID, limit)
	if err != nil {
		return polls, err
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
			return polls, err
		}
		// get time difference in human readable format
		t := utilities.TimeDiff(time)
		polls = append(polls, poll{ID: id, Title: title,
			Author: author, Time: t, NumOfVotes: numOfVotes})
	}
	err = rows.Err()
	if err != nil {
		return polls, err
	}
	// if error does not happen, return results
	return polls, nil
}

// fpQuery picks the most suitable sql query based on maxID of poll and returns
// sql rows and error
func fpQuery(maxID int, limit int) (*sql.Rows, error) {
	if maxID == 0 { // maxID = 0 when we perform the first query on "/"
		rows, err := global.DB.Query(`SELECT poll.id, poll.title,
								  users.username, poll.time,
								  (select count(*) as votes from vote where vote.poll_id = poll.id)
								  FROM poll
								  LEFT JOIN users on users.id = poll.created_by
								  ORDER BY poll.id desc
								  limit $1`, limit)
		return rows, err
	}

	rows, err := global.DB.Query(`SELECT poll.id, poll.title,
								  users.username, poll.time,
								  (select count(*) as votes from vote where vote.poll_id = poll.id)
								  FROM poll
								  LEFT JOIN users on users.id = poll.created_by
								  WHERE poll.id <= $1
								  ORDER BY poll.id desc
								  limit $2`, maxID, limit)
	return rows, err
}
