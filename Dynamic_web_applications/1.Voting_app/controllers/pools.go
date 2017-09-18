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
	Title  string
	Option []string
}

// ViewPool takes care for displaying existing pools
func ViewPool(w http.ResponseWriter, r *http.Request) {
	poolID := r.URL.Path
	poolID = strings.Split(poolID, "/")[2]

	rows, err := global.DB.Query(`SELECT title, option from pool
									LEFT JOIN poolOption
									on pool.id = poolOption.pool_id
									where pool.id = $1;`, poolID)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "This pool does not exist", http.StatusInternalServerError)
		return
	}

	var (
		title      string
		poolOption string
	)

	pool := Pool{}
	// parsing rows from database
	for rows.Next() {
		err := rows.Scan(&title, &poolOption)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		pool.Title = title
		pool.Option = append(pool.Option, poolOption)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// check if pool title exists, if it doesn't => display the 404 page
	if len(pool.Title) > 0 && len(pool.Option) > 0 {
		t := template.Must(template.ParseFiles("templates/voteDetails.html",
			"templates/navbar.html", "templates/styles.html"))
		err = t.Execute(w, pool)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
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
