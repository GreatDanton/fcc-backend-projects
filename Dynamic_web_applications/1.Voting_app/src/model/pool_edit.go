package model

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// PollUpdate handles poll updates with new data posted from poll_edit template
func PollUpdate(poll Poll, newPollTitle string, newPollOptions [][]string) error {
	// update database with new data
	pollID := poll.ID
	tx, err := global.DB.Begin()
	if err != nil {
		return err
	}

	// update poll title
	err = updatePollTitle(pollID, newPollTitle, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = deleteUnusedVoteOptions(poll.Options, newPollOptions, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update poll options -> we preserve votes of the existing options
	for _, option := range newPollOptions {
		err = updatePollOptions(pollID, option, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// if no error happened while interacting with database
	tx.Commit()
	return nil
}

// updatePollTitle updates poll title of the poll with the id of pollID
func updatePollTitle(pollID string, updatedTitle string, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE poll
							 SET title = $1
							 where poll.id = $2`)
	if err != nil {
		return fmt.Errorf("updatePollTitle: Could not prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedTitle, pollID)
	if err != nil {
		return fmt.Errorf("editPollSubmit: Could not update poll title: %v", err)
	}
	return nil // if everything is ok
}

// updatePollOptions updates poll options based on editPoll template:
// - options that are unchanged stay the same so we preserve the vote count
// - options that changed are inserted into options table
func updatePollOptions(pollID string, option []string, tx *sql.Tx) error {
	// option is in form [option-id]
	newTitle := option[0]
	id, err := strconv.Atoi(option[1])
	if err != nil {
		return err
	}
	var dbTitle string
	err = global.DB.QueryRow(`SELECT option from polloption
							   WHERE id = $1 and poll_id = $2`, id, pollID).Scan(&dbTitle)
	if err != nil {
		// empty rows, add new option to polloption table
		if err == sql.ErrNoRows {
			err = addPollOption(pollID, newTitle, tx)
			if err != nil {
				return fmt.Errorf("error while adding new poll option id=%v: %v", id, err)
			}
			return nil
		}
		// an actual error happened
		return fmt.Errorf("error while parsing polloption id=%v from db: %v", id, err)
	}
	// if ids are same and titles are different
	if newTitle != dbTitle {
		stmt, err := tx.Prepare(`UPDATE polloption
						   SET option = $1
						   WHERE id = $2`)
		if err != nil {
			return fmt.Errorf("error while preparing update statement id=%v: %v", id, err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(newTitle, id)
		if err != nil {
			return fmt.Errorf("error executing poll update statement: %v", err)
		}
	}
	// everything is allright, no error happened
	return nil
}

// getUnusedVoteOptions reuturns array of ids that should be deleted from
// the database, because the user changed them in edit poll template
func getUnusedVoteOptions(dbIDs []string, newIDs []string) []string {
	deleteIDs := []string{}
	for _, dbID := range dbIDs {
		//title := option[0]
		exist := utilities.StringInSlice(dbID, newIDs)
		if !exist {
			deleteIDs = append(deleteIDs, dbID)
		}
	}
	return deleteIDs
}

// deleteUnusedVoteOptions removes unused voteoptions from database
func deleteUnusedVoteOptions(dbOptions [][]string, newOptions [][]string, tx *sql.Tx) error {
	// getUnusedVoteOptions
	newPollOptionIDs := getOptionIDs(newOptions)
	dbPollOptionIDs := getOptionIDs(dbOptions)
	unusedIDs := getUnusedVoteOptions(dbPollOptionIDs, newPollOptionIDs)
	// delete unused IDs
	for _, id := range unusedIDs {
		_, err := tx.Exec(`delete from polloption where id = $1`, id)
		if err != nil {
			return fmt.Errorf("deleteUnusedVoteOptions: could not delete option id=%v: %v", id, err)
		}
	}
	// everything is okay
	return nil
}

// getOptionIDs returns array of ids
// input data options should look like: [[title1, id1], [title2, id2]]
func getOptionIDs(options [][]string) []string {
	IDarr := make([]string, 0, len(options))
	for _, option := range options {
		//option title := option[0]
		id := option[1]
		IDarr = append(IDarr, id)
	}
	return IDarr
}
