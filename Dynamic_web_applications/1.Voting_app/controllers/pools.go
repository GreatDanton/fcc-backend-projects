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
	Author  string
	Title   string
	Options []string
	Votes   [][]string // contains [vote Option, vote count]
}

// ViewPool takes care for displaying existing pools in /view/pool_id
func ViewPool(w http.ResponseWriter, r *http.Request) {
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

	// check if pool title exists, if it doesn't => display the 404 page
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

// getPoolDetails returns Title, Author, Vote options from database for chosen poolID
func getPoolDetails(poolID string) (Pool, error) {
	pool := Pool{}
	rows, err := global.DB.Query(`SELECT title, users.username, option from pool
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
		title      string
		author     string
		poolOption string
	)
	// parsing rows from database
	for rows.Next() {
		err := rows.Scan(&title, &author, &poolOption)
		if err != nil {
			return pool, err
		}
		pool.Title = title
		pool.Author = author
		pool.Options = append(pool.Options, poolOption)
		// add number of votes
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		return pool, err
	}

	return pool, nil
}

// get vote count for pool with chosen poolID
// returns [[Vote option 1, count 1], [Vote option 2, count 2]]
func getPoolVotes(poolID string) ([][]string, error) {
	votes := [][]string{} //Votes{}
	// returns: optionName, number of votes (sorted descending)
	rows, err := global.DB.Query(`SELECT poolOption.option, count(vote.option_id) from pooloption
								  LEFT JOIN vote
								  on pooloption.id = vote.option_id
								  where pooloption.pool_id = $1
								  group by poolOption.option
								  order by count desc`, poolID)
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

		// sorting strings in ascending order
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
		url := fmt.Sprintf("/view/%v", poolID)
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
