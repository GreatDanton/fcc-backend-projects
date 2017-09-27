package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// Pool structure used to parse values from database
type Pool struct {
	Author        string     // author username
	Title         string     // title of the pool
	Options       [][]string // contains [[option title, option_id]]
	Votes         [][]string // contains [[vote Option, vote count]]
	ErrorPostVote string     // display error when user submits his vote
	LoggedInUser  User       // User struct for rendering different templates based on user login status
}

// ViewPool takes care for handling existing pools in /pool/pool_id
// displaying existing pools and handling voting part of the pool
func ViewPool(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayPool(w, r, Pool{})
	case "POST":
		poolPostHandler(w, r)
	default:
		displayPool(w, r, Pool{})
	}
}

// displayPool is handling GET request for VIEW POOL function
// displayPool displays data for chosen pool /pool/:id and returns
//404 page if pool does not exist
func displayPool(w http.ResponseWriter, r *http.Request, poolMsg Pool) {
	poolID := strings.Split(r.URL.EscapedPath(), "/")[2]

	pool, err := getPoolDetails(poolID)
	if err != nil {
		fmt.Println("Error while getting pool details: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	votes, err := getPoolVotes(poolID)
	if err != nil {
		fmt.Printf("Error while getting pool votes count: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	pool.Votes = votes
	user := LoggedIn(r)
	pool.LoggedInUser = user

	// check if pool title exists and display relevant template with
	// pool data filled in
	if len(pool.Title) > 0 && len(pool.Options) > 0 {
		// TODO: fix this ugly implementation
		pool.ErrorPostVote = poolMsg.ErrorPostVote
		err = global.Templates.ExecuteTemplate(w, "details", pool)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else { // if db does not return any rows -> pool does not exist, display 404
		err := global.Templates.ExecuteTemplate(w, "404", pool)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// poolPostHandler handles different methods of form submit
func poolPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.Form["_method"][0]
	switch method {
	case "post":
		postVote(w, r) // user posted vote
	case "put":
		editPool(w, r) // user wanted to edit pool
	case "delete":
		deletePool(w, r) // user wanted to delete pool
	default:
		postVote(w, r)
	}
}

// postVote function handles posting votes on each /pool/:id
func postVote(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in, if it's not return 403 forbidden
	user := LoggedIn(r)
	if !user.LoggedIn {
		http.Redirect(w, r, r.URL.Path, http.StatusForbidden)
		return
	}
	r.ParseForm()
	var optionID string
	poolID := strings.Split(r.URL.EscapedPath(), "/")[2]

	for _, value := range r.Form { // f.Form contains only one option
		optionID = value[0]
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

//
// getPoolDetails returns {poolID, Title, Author, [pooloption, pooloptionID]}
// from database for chosen poolID
func getPoolDetails(poolID string) (Pool, error) {
	pool := Pool{}
	rows, err := global.DB.Query(`SELECT title, users.username, pooloption.option,
								  pooloption.id from pool
								  LEFT JOIN poolOption
								  on pool.id = poolOption.pool_id
								  LEFT JOIN users
								  on users.id = pool.created_by
								  where pool.id = $1;`, poolID)
	if err != nil {
		return pool, err
	}
	defer rows.Close()

	// defining variables for parsing rows from db
	var (
		title        string
		author       string
		poolOption   string
		poolOptionID string
	)
	// parsing rows from database
	for rows.Next() {
		err := rows.Scan(&title, &author, &poolOption, &poolOptionID)
		if err != nil {
			return pool, err
		}
		pool.Title = title
		pool.Author = author

		option := []string{poolOption, poolOptionID}
		pool.Options = append(pool.Options, option)
	}

	err = rows.Err()
	if err != nil {
		return pool, err
	}

	return pool, nil
}

// getPoolVotes returns vote count for pool with chosen poolID
// returns [[Vote option 1, count 1], [Vote option 2, count 2]]
func getPoolVotes(poolID string) ([][]string, error) {
	votes := [][]string{} //Votes{}
	// returns: optionID, optionName, number of votes => sorted by increasing id
	// this ensures vote options results are returned the same way as they were posted
	rows, err := global.DB.Query(`SELECT poolOption.id, poolOption.option,
								  count(vote.option_id) from pooloption
								  LEFT JOIN vote
								  on pooloption.id = vote.option_id
								  where pooloption.pool_id = $1
								  group by poolOption.id
								  order by poolOption.id asc`, poolID)
	if err != nil {
		return votes, err
	}
	defer rows.Close()

	var (
		id         string
		voteOption string
		count      string
	)
	for rows.Next() {
		err := rows.Scan(&id, &voteOption, &count)
		if err != nil {
			return votes, err
		}
		// appending results of table rows to Votes
		vote := []string{voteOption, count}
		votes = append(votes, vote)
	}
	err = rows.Err()
	if err != nil {
		return votes, err
	}
	return votes, nil
}

//newPoolError struct is used to display error messages in
// new pool template
type newPoolError struct {
	Title            string
	TitleError       string
	VoteOptionsError string
	LoggedInUser     User
}

// CreateNewPool takes care of handling creation of the new pool in url: /new
func CreateNewPool(w http.ResponseWriter, r *http.Request) {
	user := LoggedIn(r)
	// check if user is logged in, otherwise redirect to /login page
	if !user.LoggedIn {
		fmt.Println("")
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	errMsg := newPoolError{LoggedInUser: user}

	if r.Method == "GET" {
		err := global.Templates.ExecuteTemplate(w, "newPool", errMsg)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		poolTitle := r.Form["poolTitle"][0]
		poolTitle = strings.TrimSpace(poolTitle)
		// check if poolTitle exists else return template with error message
		if len(poolTitle) < 1 {
			e := newPoolError{TitleError: "Please add title to your pool"}
			err := global.Templates.ExecuteTemplate(w, "newPool", e)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		order := make([]string, 0, len(r.Form))
		// r.Form returns a map, we have to add fields in db in correct order
		// so we don't confuse the end user, why their options are borked
		// => that is in the same order the user wanted to post options
		for key, option := range r.Form {
			voteOption := strings.TrimSpace(option[0])     // trim empty space from pool option
			if key != "poolTitle" && len(voteOption) > 0 { // filter out empty fields and title
				order = append(order, key)
			}
		}
		// if there are not at least 2 options to vote for return error into template
		if len(order) < 2 {
			e := newPoolError{Title: poolTitle, VoteOptionsError: "Please add at least two options"}
			err := global.Templates.ExecuteTemplate(w, "newPool", e)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Internal Server error", http.StatusInternalServerError)
			}
			return
		}

		// this ensures pool options are inserted into database in
		// the same order as the end-user intended
		sort.Strings(order)
		voteOptions := make([]string, 0, len(order))
		for _, value := range order {
			voteOptions = append(voteOptions, value)
		}

		// Adding new pool into database => begin SQL transaction
		// all inserts must succeed
		tx, err := global.DB.Begin()
		if err != nil {
			fmt.Println(err)
			return
		}

		poolID, err := addPoolTitle(poolTitle, user, tx)
		if err != nil {
			fmt.Printf("addPoolTitle: %v\n", err)
			tx.Rollback()
			return
		}

		// insert posts into postOptions database
		for _, value := range voteOptions {
			option := r.Form[value][0] // text of the voteOption
			err := addPoolOption(poolID, option, tx)
			if err != nil {
				fmt.Printf("addPoolOption: %v\n", err)
				tx.Rollback()
				return
			}
		}
		// end of SQL transaction
		tx.Commit() // if no errors occur, commit to database
		// redirect to new post with status code 303
		url := fmt.Sprintf("/pool/%v", poolID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// add post title to database
func addPoolTitle(title string, user User, tx *sql.Tx) (string, error) {
	// get user id from currently logged in user
	userID := user.ID

	var id string
	err := tx.QueryRow(`INSERT into pool(created_by, title)
							 values($1, $2) RETURNING id`, userID, title).Scan(&id)
	if err != nil {
		return "", err
	}

	poolID := fmt.Sprintf("%v", id)
	return poolID, nil
}

// add new post questions to database
func addPoolOption(poolID, option string, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT into poolOption(pool_id, option)
							 values($1, $2);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(poolID, option)
	if err != nil {
		return err
	}
	return nil
}

// editPool handles edit button press on poolDetails page
func editPool(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Edit pool")
}

func deletePool(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete pool")
}
