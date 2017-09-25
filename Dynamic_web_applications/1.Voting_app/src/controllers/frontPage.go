package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

type pool struct {
	ID         string
	Time       string
	Title      string
	Author     string
	NumOfVotes string
}

// FrontPage takes care of displaying front page of Voting Application
func FrontPage(w http.ResponseWriter, r *http.Request) {
	user := User{Username: "test", ID: "id_test"}
	CreateToken(user)

	// getting database response
	rows, err := global.DB.Query(`SELECT pool.id, pool.title,
									users.username, pool.time,
									(select count(*) as votes from vote where vote.pool_id = pool.id)
									FROM pool
									LEFT JOIN users on users.id = pool.created_by
									ORDER BY pool.id desc limit 20`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id         string
		title      string
		author     string
		time       time.Time
		numOfVotes string
	)

	pools := []pool{}

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &time, &numOfVotes)
		if err != nil {
			log.Fatal(err)
		}
		// get time difference in human readable format
		t := utilities.TimeDiff(time)
		pools = append(pools, pool{ID: id, Title: title,
			Author: author, Time: t, NumOfVotes: numOfVotes})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// displaying template
	t := template.Must(template.ParseFiles("templates/index.html", "templates/navbar.html"))
	err = t.Execute(w, pools)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
