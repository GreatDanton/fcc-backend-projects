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

// EditPoolHandler handles displaying edit pool template and submitting
// updates to the database
func EditPoolHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		editPoolView(w, r)
	case "POST":
		editPoolSubmit(w, r)
	default:
		editPoolView(w, r)
	}
}

// editPool handles edit button press on poolDetails page
func editPoolView(w http.ResponseWriter, r *http.Request) {
	// insert stuff into fields
	//
	loggedUser := LoggedIn(r)
	if !loggedUser.LoggedIn {
		err := global.Templates.ExecuteTemplate(w, "403", http.StatusForbidden)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	poolID := strings.Split(r.URL.Path, "/")[2] //ustrings.Split(u, "/")[2]
	pool, err := getPoolDetails(poolID)
	if err != nil {
		fmt.Println(err)
		return
	}
	pool.LoggedInUser = loggedUser
	if loggedUser.Username != pool.Author {
		fmt.Println("Currently logged in user is not the author")
		err := global.Templates.ExecuteTemplate(w, "403", pool)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = global.Templates.ExecuteTemplate(w, "editPool", pool)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

//editPoolSubmit handles pool title and pool options updates
func editPoolSubmit(w http.ResponseWriter, r *http.Request) {
	// get pool id
	poolID := mux.Vars(r)["poolID"]
	user := LoggedIn(r)

	// get title, options from template
	updatedTitle, optionFieldNames, err := parsePoolParams(w, r, "edit")
	if err != nil {
		// error template rendering is already done in parsePoolParams
		fmt.Println("parsePoolParams:", err)
		return
	}

	pool, err := getPoolDetails(poolID)

	// check if logged in user is pool author
	if user.Username != pool.Author {
		fmt.Println("editPoolSubmit: currently logged in user is not the author of the pool")
		err := global.Templates.ExecuteTemplate(w, "403", nil)
		if err != nil {
			fmt.Println("editPoolSubmit: ExecuteTemplate:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	// update database with new data
	tx, err := global.DB.Begin()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error, could not update your pool", http.StatusInternalServerError)
		return
	}

	// update pool title
	err = updatePoolTitle(poolID, updatedTitle, tx)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		http.Error(w, "Internal server error, pool could not be updated", http.StatusInternalServerError)
		return
	}

	newPoolOptions := [][]string{}
	for _, option := range optionFieldNames {
		// option looks like [option-1]
		optionTitle := r.Form[option][0]
		id := strings.Split(option, "-")[1]
		arr := []string{optionTitle, id}
		newPoolOptions = append(newPoolOptions, arr)
	}

	err = deleteUnusedVoteOptions(pool.Options, newPoolOptions, tx)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return
	}

	// update pool options -> we preserve votes of the existing options
	for _, option := range newPoolOptions {
		err = updatePoolOptions(poolID, option, tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
	}
	// if no error happened while interacting with database
	tx.Commit()
	url := fmt.Sprintf("/pool/%v", poolID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// updatePoolTitle updates pool title of the pool with the id of poolID
func updatePoolTitle(poolID string, updatedTitle string, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE pool
							 SET title = $1
							 where pool.id = $2`)
	if err != nil {
		return fmt.Errorf("updatePoolTitle: Could not prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedTitle, poolID)
	if err != nil {
		return fmt.Errorf("editPoolSubmit: Could not update pool title: %v", err)
	}
	return nil // if everything is ok
}

// updatePoolOptions updates pool options based on editPool template:
// - options that are unchanged stay the same so we preserve the vote count
// - options that changed are inserted into options table
func updatePoolOptions(poolID string, option []string, tx *sql.Tx) error {
	// option is in form [option-id]
	newTitle := option[0]
	id, err := strconv.Atoi(option[1])
	if err != nil {
		return err
	}
	var dbTitle string
	err = global.DB.QueryRow(`SELECT option from pooloption
							   WHERE id = $1 and pool_id = $2`, id, poolID).Scan(&dbTitle)
	if err != nil {
		// empty rows, add new option to pooloption table
		if err == sql.ErrNoRows {
			err = addPoolOption(poolID, newTitle, tx)
			if err != nil {
				return fmt.Errorf("error while adding new pool option id=%v: %v", id, err)
			}
			return nil
		}
		// an actual error happened
		return fmt.Errorf("error while parsing pooloption id=%v from db: %v",
			id, err)
	}

	// if ids are same and titles are different
	if newTitle != dbTitle {
		stmt, err := tx.Prepare(`UPDATE pooloption
						   SET option = $1
						   WHERE id = $2`)
		if err != nil {
			return fmt.Errorf("error while preparing update statement id=%v: %v", id, err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(newTitle, id)
		if err != nil {
			return fmt.Errorf("error executing pool update statement: %v", err)
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
	newPoolOptionIDs := getOptionIDs(newOptions)
	dbPoolOptionIDs := getOptionIDs(dbOptions)
	unusedIDs, err := getUnusedVoteOptions(dbPoolOptionIDs, newPoolOptionIDs)
	if err != nil {
		return err
	}
	// delete unused IDs
	for _, id := range unusedIDs {
		_, err = tx.Exec(`delete from pooloption where id = $1`, id)
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
// parsePoolParams fetches data from editPool/newPool templates form and returns:
// poolTitle, [voteOptions], error
func parsePoolParams(w http.ResponseWriter, r *http.Request, template string) (string, []string, error) {
	if template == "edit" {
		template = "editPool"
	} else {
		template = "newPool"
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", []string{}, err
	}
	errMsg := Pool{}

	poolTitle := strings.TrimSpace(r.Form["poolTitle"][0])
	// check if poolTitle exists else return template with error message
	if len(poolTitle) < 1 {
		errMsg.Errors.TitleError = "Please add title to your pool"
		//e := newPoolError{TitleError: "Please add title to your pool"}
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
		voteOption := strings.TrimSpace(option[0])     // trim empty space from pool option
		if key != "poolTitle" && len(voteOption) > 0 { // filter out empty fields and title
			order = append(order, key)
		}
	}
	// if there are not at least 2 options to vote for return error into template
	if len(order) < 2 {
		errMsg.Errors.Title = poolTitle
		errMsg.Errors.VoteOptionsError = "Please add at least two options"
		// add vote options to the pool struct, otherwise options are missing upon
		// template rerender
		errMsg.Errors.VoteOptions = order
		//e := newPoolError{Title: poolTitle, VoteOptionsError: "Please add at least two options"}
		err := global.Templates.ExecuteTemplate(w, template, errMsg)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return "", []string{}, err
		}
		return "", []string{}, fmt.Errorf("User added less than 2 vote options")
	}

	// this ensures pool options are inserted into database in
	// the same order as the end-user intended
	sort.Strings(order)
	voteOptions := make([]string, 0, len(order))
	for _, value := range order {
		voteOptions = append(voteOptions, value)
	}

	return poolTitle, voteOptions, nil
}
