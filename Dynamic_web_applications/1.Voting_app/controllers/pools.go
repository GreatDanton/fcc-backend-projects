package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/global"
)

// Pool structure used to parse values from database
type Pool struct {
	ID      string
	Author  string
	Title   string
	Options [][]string // contains [[option title, option_id]]
	Votes   [][]string // contains [vote Option, vote count]
}

// ViewPool takes care for handling existing pools in /pool/pool_id
// displaying existing pools and handling voting part of the pool
func ViewPool(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		displayPool(w, r)
	case "POST":
		postVote(w, r)
	default:
		displayPool(w, r)
	}
}

// displayPool is handling GET request for VIEW POOL function
// displayPool displays data for chosen pool /pool/:id and returns
//404 page if pool does not exist
func displayPool(w http.ResponseWriter, r *http.Request) {
	poolID := r.URL.Path
	poolID = strings.Split(poolID, "/")[2]

	pool, err := getPoolDetails(poolID)
	if err != nil {
		fmt.Println("Error while getting pool details: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	votes, err := getPoolVotes(poolID)
	if err != nil {
		fmt.Printf("Error while getting pool votes count: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	pool.Votes = votes

	// check if pool title exists and display relevant template with
	// pool data filled in
	if len(pool.Title) > 0 && len(pool.Options) > 0 {
		t := template.Must(template.ParseFiles("templates/voteDetails.html",
			"templates/navbar.html", "templates/styles.html"))
		err = t.Execute(w, pool)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else { // if db does not return any rows -> pool does not exist, display 404
		t := template.Must(template.ParseFiles("templates/404.html",
			"templates/navbar.html",
			"templates/styles.html"))
		err := t.Execute(w, "")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// postVote takes care of POST request on ViewPOOL
// This function handles posting votes on each /pool/:id
func postVote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var optionID string
	poolID := r.Form["poolID"][0]

	for key, value := range r.Form {
		if key != "poolID" {
			optionID = value[0]
		}
	}

	// userID should be logged in user -> authentication part is still missing
	// voting as User1 for now
	userID := 1
	// check if vote for user already exists

	// add vote to database
	_, err := global.DB.Exec(`INSERT into vote(pool_id, option_id, voted_by)
	values($1, $2, $3)`, poolID, optionID, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// refresh page -> redirect to the same page
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

// getPoolDetails returns Title, Author, Vote options from database for chosen poolID
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

	// defining variables for parsing rows from db
	var (
		title        string
		author       string
		poolOption   string
		poolOptionID string
	)
	// assign ID
	pool.ID = poolID
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
		// add number of votes
	}
	defer rows.Close()

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
	// returns: optionName, number of votes (sorted descending)
	rows, err := global.DB.Query(`SELECT poolOption.option, count(vote.option_id) from pooloption
								  LEFT JOIN vote
								  on pooloption.id = vote.option_id
								  where pooloption.pool_id = $1
								  group by poolOption.option`, poolID)
	if err != nil {
		return votes, err
	}
	defer rows.Close()

	var (
		voteOption string
		count      string
	)
	for rows.Next() {
		err := rows.Scan(&voteOption, &count)
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
}

// CreateNewPool takes care of handling creation of the new pool in url: /new
func CreateNewPool(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/newPool.html", "templates/navbar.html", "templates/styles.html"))
	if r.Method == "GET" {
		err := t.Execute(w, nil)
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
			err := t.Execute(w, e)
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
			err := t.Execute(w, e)
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

		poolID, err := addPoolTitle(poolTitle, tx)
		if err != nil {
			fmt.Printf("addPoolTitle: %v", err)
			tx.Rollback()
			return
		}

		// insert posts into postOptions database
		for _, value := range voteOptions {
			option := r.Form[value][0] // text of the voteOption
			err := addPoolOption(poolID, option, tx)
			if err != nil {
				fmt.Printf("addPoolOption: %v", err)
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
func addPoolTitle(title string, tx *sql.Tx) (string, error) {
	// get user id from currently logged in user
	userID := 1 // for now using 1

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
