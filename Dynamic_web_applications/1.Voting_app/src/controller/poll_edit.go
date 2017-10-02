package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

// EditPollHandler handles displaying edit poll template and submitting
// updates to the database
func EditPollHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		editPollView(w, r)
	case "POST":
		editPollSubmit(w, r)
	default:
		editPollView(w, r)
	}
}

// editPollView handles edit button press on pollDetails page
func editPollView(w http.ResponseWriter, r *http.Request) {
	// insert stuff into fields
	loggedUser := model.LoggedInUser(r)
	if !loggedUser.LoggedIn {
		err := templates.Execute(w, "403", nil)
		if err != nil {
			fmt.Println("editPollView:", err)
		}
		return
	}

	pollID := strings.Split(r.URL.Path, "/")[2] //ustrings.Split(u, "/")[2]
	poll, err := model.GetPollDetails(pollID)
	if err != nil {
		fmt.Println(err)
		return
	}
	poll.LoggedInUser = loggedUser
	if loggedUser.Username != poll.Author {
		fmt.Println("Currently logged in user is not the author")
		err := templates.Execute(w, "403", poll)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	err = global.Templates.ExecuteTemplate(w, "editPoll", poll)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//editPollSubmit handles poll title and poll options updates =>
// post request on edit
func editPollSubmit(w http.ResponseWriter, r *http.Request) {
	pollID := bone.GetValue(r, "pollID")
	user := model.LoggedInUser(r)

	// get title, options from template
	newPollTitle, optionFieldNames, err := parsePollParams(w, r, "edit")
	if err != nil {
		// error template rendering is already done in parsePollParams
		fmt.Println("parsePollParams:", err)
		return
	}

	poll, err := model.GetPollDetails(pollID)
	// check if logged in user is poll author
	if user.Username != poll.Author {
		fmt.Println("editPollSubmit: currently logged in user is not the author of the poll")
		err = templates.Execute(w, "403", nil)
		if err != nil {
			fmt.Println("editPollSubmit: ExecuteTemplate:", err)
		}
		return
	}

	// parse new poll options from edit template input_fields
	// newPollOptions = [[optionTitle1, id1], [optionTitle2, id2]]
	newPollOptions := [][]string{}
	for _, option := range optionFieldNames {
		// option looks like [option-1]
		optionTitle := r.Form[option][0]
		id := strings.Split(option, "-")[1]
		arr := []string{optionTitle, id}
		newPollOptions = append(newPollOptions, arr)
	}
	err = model.PollUpdate(poll, newPollTitle, newPollOptions)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("/poll/%v", pollID)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
