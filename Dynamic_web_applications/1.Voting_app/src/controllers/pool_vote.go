package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

// postVote function handles posting votes on each /pool/:id
func postVote(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in, if it's not return 403 forbidden
	user := LoggedIn(r)
	if !user.LoggedIn {
		http.Redirect(w, r, r.URL.Path, http.StatusForbidden)
		return
	}
	r.ParseForm()
	poolID := strings.Split(r.URL.EscapedPath(), "/")[2]
	// get optionID, if the user did not pick anything
	// optionID is empty string
	var optionID string
	for key, value := range r.Form {
		if key == "voteOption" {
			optionID = value[0]
		}
	}

	poolMsg := Pool{}
	// if no vote option was chosen rerender template and display
	// error message to user
	if optionID == "" {
		poolMsg.ErrorPostVote = "Please pick your vote option"
		fmt.Println("postVote: no vote option was chosen")
		displayPool(w, r, poolMsg)
		return
	}

	// check if user is changing vote options via html, this prevents
	// spamming votes for options that do not exist for this poolID
	voteOptions, err := getVoteOptions(poolID)
	if err != nil {
		fmt.Println("postVote:", "getVoteOptions:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("OptionID:", optionID)
	fmt.Println("voteOPTIONS:", voteOptions)
	ok := utilities.StringInSlice(optionID, voteOptions)
	if !ok {
		poolMsg.ErrorPostVote = "You'll have to be more clever."
		fmt.Println("PostVote:", "User is changing vote options")
		displayPool(w, r, poolMsg)
		return
	}

	// use user id of logged in user
	userID := user.ID

	// check if vote for user already exists
	var dbVoteID string
	var dbOption string
	err = global.DB.QueryRow(`SELECT id, option_id from vote
							   WHERE voted_by = $1
							   AND pool_id = $2`, userID, poolID).Scan(&dbVoteID, &dbOption)

	if err != nil {
		// if user did not vote, add users vote into database
		if err == sql.ErrNoRows {
			// add vote to database
			_, e := global.DB.Exec(`INSERT into vote(pool_id, option_id, voted_by)
									  values($1, $2, $3)`, poolID, optionID, userID)
			if e != nil {
				fmt.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			// refresh page -> redirect to the same page
			// and prints http: multiple response.WriteHeader calls
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}

		// if an actual error occured, display internal server error msg
		fmt.Printf("postVote: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
			fmt.Printf("postVote: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	// refresh page upon successful POST request
	//(does not matter if db was updated => we should always redirect after POST)
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

// getVoteOptions returns array of vote options that exist in db
// for chosen poolID.
func getVoteOptions(poolID string) ([]string, error) {
	options := []string{}
	rows, err := global.DB.Query(`SELECT id from pooloption
								  WHERE pool_id = $1`, poolID)
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
