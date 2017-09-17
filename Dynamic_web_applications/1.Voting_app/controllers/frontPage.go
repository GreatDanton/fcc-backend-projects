package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/global"
)

type question struct {
	ID       string
	Time     string
	Question string
}

// FrontPage takes care of displaying front page of Voting Application
func FrontPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// getting database response
		rows, err := global.DB.Query("SELECT * from QUESTIONS")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var (
			id    string
			time  string
			quest string
		)

		Questions := []question{}

		for rows.Next() {
			err := rows.Scan(&id, &time, &quest)
			if err != nil {
				log.Fatal(err)
			}
			Questions = append(Questions, question{ID: id, Time: time, Question: quest})
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		// displaying template
		t := template.Must(template.ParseFiles("templates/index.html"))
		err = t.Execute(w, Questions)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		fmt.Println("Posting stuff, handle with db")
	}

}
