package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

//newPoolError struct is used to display error messages in
// new pool template
type newPoolError struct {
	Title            string
	TitleError       string
	VoteOptionsError string
	LoggedInUser     User
}

// CreateNewPool takes care of handling creation of the new pool in url: /new
// add post title to database
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
		poolTitle, voteOptions, err := parsePoolParams(w, r)
		if err != nil {
			// displaying error message is done in function
			fmt.Println(err)
			return
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
