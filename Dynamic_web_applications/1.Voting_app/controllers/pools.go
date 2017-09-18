package controllers

import (
	"fmt"
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
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	if len(pool.Title) > 0 {
		fmt.Fprintf(w, fmt.Sprintf("%v", pool))
	} else {
		fmt.Fprintf(w, "This pool does not exist")
	}
}
