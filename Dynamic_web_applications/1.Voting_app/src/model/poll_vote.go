package model

import (
	"database/sql"
	"fmt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// GetVoteOptions returns array of vote options that exist in db
// for chosen poll with id = pollID.
func GetVoteOptions(pollID string) ([]string, error) {
	options := []string{}
	rows, err := global.DB.Query(`SELECT id from polloption
								  WHERE poll_id = $1`, pollID)
	if err != nil {
		return options, err
	}
	defer rows.Close()

	var (
		id string
	)
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return options, err
		}
		options = append(options, id)
	}

	err = rows.Err()
	if err != nil {
		return options, err
	}
	return options, nil
}

// PollAddUserVote checks if user already voted on the chosen poll
func PollAddUserVote(pollID, optionID, userID string) error {
	// check if vote for user already exists
	var dbVoteID string
	var dbOption string
	err := global.DB.QueryRow(`SELECT id, option_id from vote
							   WHERE voted_by = $1
							   AND poll_id = $2`, userID, pollID).Scan(&dbVoteID, &dbOption)

	if err != nil {
		// if user did not vote, add users vote into database
		if err == sql.ErrNoRows {
			// add vote to database
			_, e := global.DB.Exec(`INSERT into vote(poll_id, option_id, voted_by)
									  values($1, $2, $3)`, pollID, optionID, userID)
			if e != nil {
				fmt.Println("PollAddUserVote: Error inserting into vote database:", err)
				return err
			}
			// Inserting new vote into database was successful
			// return function and redirect in the outer function
			return nil
		}

		// if an actual error occured, display internal server error msg
		fmt.Println("PollAddUserVote:", err)
		return err
	}

	// error did not occur user already voted -> change his vote
	// if his recent vote option is different than his past
	// update database with his recent option
	if optionID != dbOption {
		// if user change his mind, update his vote
		_, err = global.DB.Exec(`UPDATE vote SET
								 option_id = $1
								 where id = $2`, optionID, dbVoteID)
		if err != nil {
			fmt.Println("PollAddUserVote: updating vote error:", err)
			return err
		}
	}
	// error did not occur, return nil and refresh page in the outer function
	return nil
}
