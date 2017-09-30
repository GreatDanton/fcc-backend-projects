package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

// EditPollHandler handles displaying edit poll template and submitting
// updates to the database
func EditPollHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		editPollView(w, r)
	case "POST":
		editPollSubmit(w, r)
	default:
		editPollView(w, r)
	}
}

// editPoll handles edit button press on pollDetails page
func editPollView(w http.ResponseWriter, r *http.Request) {
	// insert stuff into fields
	//
	loggedUser := LoggedIn(r)
	if !loggedUser.LoggedIn {
		err := global.Templates.ExecuteTemplate(w, "403", nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	pollID := strings.Split(r.URL.Path, "/")[2] //ustrings.Split(u, "/")[2]
	poll, err := getPollDetails(pollID)
	if err != nil {
		fmt.Println(err)
		return
	}
	poll.LoggedInUser = loggedUser
	if loggedUser.Username != poll.Author {
		fmt.Println("Currently logged in user is not the author")
		err := global.Templates.ExecuteTemplate(w, "403", poll)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = global.Templates.ExecuteTemplate(w, "editPoll", poll)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

//editPollSubmit handles poll title and poll options updates
func editPollSubmit(w http.ResponseWriter, r *http.Request) {
	// get poll id
	pollID := mux.Vars(r)["pollID"]
	user := LoggedIn(r)

	// get title, options from template
	updatedTitle, optionFieldNames, err := parsePollParams(w, r, "edit")
	if err != nil {
		// error template rendering is already done in parsePollParams
		fmt.Println("parsePollParams:", err)
		return
	}

	poll, err := getPollDetails(pollID)

	// check if logged in user is poll author
	if user.Username != poll.Author {
		fmt.Println("editPollSubmit: currently logged in user is not the author of the poll")
		err := global.Templates.ExecuteTemplate(w, "403", nil)
		if err != nil {
			fmt.Println("editPollSubmit: ExecuteTemplate:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	// update database with new data
	tx, err := global.DB.Begin()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error, could not update your poll", http.StatusInternalServerError)
		return
	}

	// update poll title
	err = updatePollTitle(pollID, updatedTitle, tx)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		http.Error(w, "Internal server error, poll could not be updated", http.StatusInternalServerError)
		return
	}

	newPollOptions := [][]string{}
	for _, option := range optionFieldNames {
		// option looks like [option-1]
		optionTitle := r.Form[option][0]
		id := strings.Split(option, "-")[1]
		arr := []string{optionTitle, id}
		newPollOptions = append(newPollOptions, arr)
	}

	err = deleteUnusedVoteOptions(poll.Options, newPollOptions, tx)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return
	}

	// update poll options -> we preserve votes of the existing options
	for _, option := range newPollOptions {
		err = updatePollOptions(pollID, option, tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
	}
	// if no error happened while interacting with database
	tx.Commit()
	url := fmt.Sprintf("/poll/%v", pollID)
	http.Redirect(w, r, url, http.StatusSeeOther)
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
		return fmt.Errorf("error while parsing polloption id=%v from db: %v",
			id, err)
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

//
func getUnusedVoteOptions(dbIDs []string, newIDs []string) ([]string, error) {
	deleteIDs := []string{}
	for _, dbID := range dbIDs {
		//title := option[0]
		exist := utilities.StringInSlice(dbID, newIDs)
		if !exist {
			deleteIDs = append(deleteIDs, dbID)
		}
	}
	return deleteIDs, nil
}

func deleteUnusedVoteOptions(dbOptions [][]string, newOptions [][]string, tx *sql.Tx) error {
	// getUnusedVoteOptions
	newPollOptionIDs := getOptionIDs(newOptions)
	dbPollOptionIDs := getOptionIDs(dbOptions)
	unusedIDs, err := getUnusedVoteOptions(dbPollOptionIDs, newPollOptionIDs)
	if err != nil {
		return err
	}
	// delete unused IDs
	for _, id := range unusedIDs {
		_, err = tx.Exec(`delete from polloption where id = $1`, id)
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

//
// parsePollParams fetches data from editPoll/newPoll templates form and returns:
// pollTitle, [voteOptions], error
func parsePollParams(w http.ResponseWriter, r *http.Request, template string) (string, []string, error) {
	if template == "edit" {
		template = "editPoll"
	} else {
		template = "newPoll"
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", []string{}, err
	}
	errMsg := Poll{}

	pollTitle := strings.TrimSpace(r.Form["pollTitle"][0])
	// check if pollTitle exists else return template with error message
	if len(pollTitle) < 1 {
		errMsg.Errors.TitleError = "Please add title to your poll"
		//e := newPollError{TitleError: "Please add title to your poll"}
		err := global.Templates.ExecuteTemplate(w, template, errMsg)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return "", []string{}, err
		}
		return "", []string{}, fmt.Errorf("Title of the post is missing")
	}

	order := make([]string, 0, len(r.Form))
	// r.Form returns a map, we have to add fields in db in correct order
	//  (=> that is in the same order the user wanted to post options)
	// so we don't confuse the end user, why their options are borked
	for key, option := range r.Form {
		voteOption := strings.TrimSpace(option[0])     // trim empty space from poll option
		if key != "pollTitle" && len(voteOption) > 0 { // filter out empty fields and title
			order = append(order, key)
		}
	}
	// if there are not at least 2 options to vote for return error into template
	if len(order) < 2 {
		errMsg.Errors.Title = pollTitle
		errMsg.Errors.VoteOptionsError = "Please add at least two options"
		// add vote options to the poll struct, otherwise options are missing upon
		// template rerender
		errMsg.Errors.VoteOptions = order
		//e := newPollError{Title: pollTitle, VoteOptionsError: "Please add at least two options"}
		err := global.Templates.ExecuteTemplate(w, template, errMsg)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return "", []string{}, err
		}
		return "", []string{}, fmt.Errorf("User added less than 2 vote options")
	}

	// this ensures poll options are inserted into database in
	// the same order as the end-user intended
	sort.Strings(order)
	voteOptions := make([]string, 0, len(order))
	for _, value := range order {
		voteOptions = append(voteOptions, value)
	}

	return pollTitle, voteOptions, nil
}
