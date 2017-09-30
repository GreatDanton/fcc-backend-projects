package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

// UserDetails is displaying details of chosen user
// details are: username and created polls
func UserDetails(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		userDetailsGET(w, r)
	default:
		userDetailsGET(w, r)
	}
}

// User is used to display user details in /u/username
type userDetails struct {
	Username     string
	Polls        []poll
	LoggedInUser User
	Pagination   pagination
}

// userDetailsGet renders userDetail template and displays users data
// username and created polls
func userDetailsGET(w http.ResponseWriter, r *http.Request) {
	user := userDetails{}
	urlUser, err := url.PathUnescape(r.URL.EscapedPath())
	if err != nil {
		fmt.Println("UserDetailsGet: Cannot unescape url")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// splits on ? and trims right /, this ensures user is able to use / in username
	urlUser = strings.TrimRight(strings.Split(urlUser, "?")[0], "/")
	user.Username = strings.TrimSpace(urlUser[len("/u/"):])
	//TODO: old way of parsing username//user.Username = strings.TrimSpace(strings.Split(urlUser, "/")[2])

	user.LoggedInUser = LoggedIn(r)

	exist, err := userExistCheck(user.Username)
	if err != nil { // user does not exist
		fmt.Println("userExistCheck:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// if user does not exist, display 404 page
	if !exist {
		fmt.Println("User does not exist")
		err = global.Templates.ExecuteTemplate(w, "404", nil)
		if err != nil {
			fmt.Println("err")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	limit := 20
	maxID := getMaxIDParam(r)
	// get polls from user
	userPolls, err := getUserPolls(user.Username, maxID, limit)
	if err != nil {
		fmt.Printf("getUserPoll: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Polls = userPolls
	p := handlePollPagination(r, maxID, userPolls, limit)
	user.Pagination = p

	err = global.Templates.ExecuteTemplate(w, "users", user)
	if err != nil {
		fmt.Printf("userDetailsGET: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// getUserPolls fetches the database and returns user polls based on the
// the function arguments
// userName: user useranem
// maxID: maximum poll id
// limit: number of polls
func getUserPolls(username string, maxID int, limit int) ([]poll, error) {
	polls := []poll{}
	var (
		id     string
		title  string
		author string
		time   time.Time
		votes  string
	)

	var rows *sql.Rows
	var err error
	if maxID == 0 {
		rows, err = global.DB.Query(`SELECT poll.id, poll.title, poll.created_by, poll.time, count(vote.poll_id)
							   		 from poll
							   		 LEFT JOIN users
							   		 on users.id = poll.created_by
							   		 LEFT JOIN vote
							   		 on vote.poll_id = poll.id
							   		 WHERE users.username = $1
							   		 GROUP BY poll.id
							   		 order by poll.id desc
							   		 limit $2`, username, limit)
	} else {
		// poll id can't be < 0
		rows, err = global.DB.Query(`SELECT poll.id, poll.title, poll.created_by, poll.time, count(vote.poll_id)
									 from poll
									 LEFT JOIN users
									 on users.id = poll.created_by
									 LEFT JOIN vote
									 on vote.poll_id = poll.id
									 WHERE users.username = $1
									 AND poll.id <= $2
									 GROUP BY poll.id
									 order by poll.id desc
									 limit $3`, username, maxID, limit)

	}

	if err != nil {
		return polls, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &time, &votes)
		if err != nil {
			return polls, err
		}
		timeAgo := utilities.TimeDiff(time) // create submitted ...ago string
		polls = append(polls, poll{ID: id, Title: title, Author: author, Time: timeAgo, NumOfVotes: votes})
	}
	err = rows.Err()
	if err != nil {
		return polls, err
	}

	return polls, nil
}
