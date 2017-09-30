package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

//newPollError struct is used to display error messages in
// new poll template
type newPollError struct {
	Title            string
	TitleError       string
	VoteOptionsError string
	LoggedInUser     User
}

// CreateNewPoll takes care of handling creation of the new poll in url: /new
// add post title to database
func CreateNewPoll(w http.ResponseWriter, r *http.Request) {
	user := LoggedIn(r)
	// check if user is logged in, otherwise redirect to /login page
	if !user.LoggedIn {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	poll := Poll{LoggedInUser: user}
	//errMsg := newPollError{LoggedInUser: user}

	if r.Method == "GET" {
		err := global.Templates.ExecuteTemplate(w, "newPoll", poll)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		pollTitle, voteOptions, err := parsePollParams(w, r, "newPoll")
		if err != nil {
			// displaying error message is done in function
			fmt.Println("CreateNewPoll:", "parsePollParams:", err)
			return
		}

		// Adding new poll into database => begin SQL transaction
		// all inserts must succeed
		tx, err := global.DB.Begin()
		if err != nil {
			fmt.Println(err)
			return
		}

		pollID, err := addPollTitle(pollTitle, user, tx)
		if err != nil {
			fmt.Printf("addPollTitle: %v\n", err)
			tx.Rollback()
			return
		}

		// insert posts into postOptions database
		for _, value := range voteOptions {
			option := r.Form[value][0] // text of the voteOption
			err := addPollOption(pollID, option, tx)
			if err != nil {
				fmt.Printf("addPollOption: %v\n", err)
				tx.Rollback()
				return
			}
		}
		// end of SQL transaction
		tx.Commit() // if no errors occur, commit to database
		// redirect to new post with status code 303
		url := fmt.Sprintf("/poll/%v", pollID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func addPollTitle(title string, user User, tx *sql.Tx) (string, error) {
	// get user id from currently logged in user
	userID := user.ID

	var id string
	err := tx.QueryRow(`INSERT into poll(created_by, title)
							 values($1, $2) RETURNING id`, userID, title).Scan(&id)
	if err != nil {
		return "", err
	}

	pollID := fmt.Sprintf("%v", id)
	return pollID, nil
}

// add new post questions to database
func addPollOption(pollID, option string, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT into pollOption(poll_id, option)
							 values($1, $2);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pollID, option)
	if err != nil {
		return err
	}
	return nil
}
