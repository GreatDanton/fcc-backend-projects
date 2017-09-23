package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/utilities"
)

type question struct {
	ID     string
	Time   string
	Title  string
	Author string
}

// FrontPage takes care of displaying front page of Voting Application
func FrontPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// getting database response
		rows, err := global.DB.Query(`SELECT pool.id, pool.title,
									users.username, pool.time
									FROM pool
									LEFT JOIN users on users.id = pool.created_by
									ORDER BY pool.id desc limit 20`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var (
			id     string
			title  string
			author string
			time   time.Time
		)

		Questions := []question{}

		for rows.Next() {
			err := rows.Scan(&id, &title, &author, &time)
			if err != nil {
				log.Fatal(err)
			}
			// get time difference in human readable format
			t := utilities.TimeDiff(time)
			Questions = append(Questions, question{ID: id, Title: title, Author: author, Time: t})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		// displaying template
		t := template.Must(template.ParseFiles("templates/index.html", "templates/navbar.html"))
		err = t.Execute(w, Questions)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

}
