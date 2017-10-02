package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

func deletePoll(w http.ResponseWriter, r *http.Request) {
	user := model.LoggedInUser(r)
	// if user is not logged in render 403 template
	if !user.LoggedIn {
		fmt.Println("User is not logged in")
		err := templates.Execute(w, "403", nil)
		if err != nil {
			fmt.Println("delete Poll: Problem parsing templates:", err)
		}
		return
	}
	pollID := utilities.GetURLSuffix(r)
	poll, err := model.GetPollDetails(pollID)
	if err != nil {
		fmt.Println("Delete: getPollDetails", err)
		return
	}

	// if title is empty, poll does not exist
	if poll.Title == "" {
		// using post -> redirect, or post -> render?
		//http.Redirect(w, r, r.URL.Path, http.StatusNotFound)
		err = templates.Execute(w, "404", nil)
		if err != nil {
			fmt.Println("deletePoll: problem parsing 404 template", err)
			return
		}
		return
	}

	// if logged in user is not author of the post return 403 Forbidden
	if user.Username != poll.Author {
		fmt.Println("Currently logged in user is not poll author")
		err := templates.Execute(w, "403", nil)
		if err != nil {
			fmt.Println("deletePoll: problem executing template: ", err)
		}
		return
	}

	// everything is allright  delete poll with pollid
	err = model.DeletePoll(poll.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	infoMsg := info{LoggedInUser: user}
	// delete was successful, inform user
	err = templates.Execute(w, "info", infoMsg)
	if err != nil {
		fmt.Println("DeletePoll: template error:", err)
		return
	}
}
