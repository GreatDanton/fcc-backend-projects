package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

type frontPage struct {
	Polls        []model.Poll
	LoggedInUser model.User
	/* 	Pagination   model.pagination */
}

// FrontPage is used to display first page of the web application
func FrontPage(w http.ResponseWriter, r *http.Request) {
	maxID := 0
	limit := 20
	polls, err := model.GetFrontPageData(maxID, limit)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user := model.LoggedInUser(r)
	fp := frontPage{Polls: polls, LoggedInUser: user}
	err = templates.Execute(w, "frontPage", fp)
	if err != nil {
		return
	}
}
