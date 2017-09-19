package controllers

import (
	"fmt"
	"html/template"
	"net/http"
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

// ViewPool takes care for displaying existing pools
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

// CreateNewPool takes care of handling creation of the new pool
func CreateNewPool(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/newPool.html", "templates/navbar.html", "templates/styles.html"))
		err := t.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		fmt.Println("posted stuff")
	}
}
