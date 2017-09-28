package controllers

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// Pool structure used to parse values from database
type Pool struct {
	ID            string
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
		path := utilities.GetURLSuffix(r)
		if path == "edit" {
			editPool(w, r)
		} else {
			displayPool(w, r, Pool{})
		}
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
	poolID := utilities.GetURLSuffix(r)

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
	case "edit":
		editPool(w, r) // user wanted to edit pool
	case "delete":
		deletePool(w, r) // user wanted to delete pool
	default:
		postVote(w, r)
	}
}

//
// getPoolDetails returns {poolID, Title, Author, [pooloption, pooloptionID]}
// from database for chosen poolID
func getPoolDetails(poolID string) (Pool, error) {
	pool := Pool{}
	rows, err := global.DB.Query(`SELECT pool.id, pool.title, users.username, pooloption.option,
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
		id           string
		title        string
		author       string
		poolOption   string
		poolOptionID string
	)
	// parsing rows from database
	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &poolOption, &poolOptionID)
		if err != nil {
			return pool, err
		}
		pool.ID = id
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
