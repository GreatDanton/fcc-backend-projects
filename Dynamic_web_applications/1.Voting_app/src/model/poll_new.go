package model

import (
	"database/sql"
	"fmt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// AddNewPoll adds new poll to the database with pollTitle, pollOptions, Author
// returns: pollID, error
func AddNewPoll(pollTitle string, pollOptions []string, user User) (string, error) {
	// Adding new poll into database => begin SQL transaction
	// all inserts must succeed
	tx, err := global.DB.Begin()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	pollID, err := addPollTitle(pollTitle, user, tx)
	if err != nil {
		fmt.Printf("addPollTitle: %v\n", err)
		tx.Rollback()
		return pollID, err
	}

	// insert posts into postOptions database
	for _, option := range pollOptions {
		err := addPollOption(pollID, option, tx)
		if err != nil {
			fmt.Printf("addPollOption: %v\n", err)
			tx.Rollback()
			return pollID, err
		}
	}
	// end of SQL transaction
	tx.Commit() // if no errors occur, commit to database
	// redirect to new post with status code 303
	return pollID, nil
}

// addPollTitle is part of the AddNewPoll transaction, that adds poll title to the database
// returns: pollID, error
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
