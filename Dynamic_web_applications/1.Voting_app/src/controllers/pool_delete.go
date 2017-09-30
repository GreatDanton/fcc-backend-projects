package controllers

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// deletePoll handles deleting chosen polls
func deletePoll(w http.ResponseWriter, r *http.Request) {
	user := LoggedIn(r)
	// if user is not logged in render 403 template
	infoMsg := info{LoggedInUser: user}

	if !user.LoggedIn {
		fmt.Println("User is not logged in")
		err := global.Templates.ExecuteTemplate(w, "403", infoMsg)
		if err != nil {
			fmt.Println("delete Poll: Problem parsing templates:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	pollID := utilities.GetURLSuffix(r)
	poll, err := getPollDetails(pollID)
	if err != nil {
		fmt.Println("Delete: getPollDetails", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// if title is empty, poll does not exist
	if poll.Title == "" {
		// using post -> redirect, or post -> render?
		//http.Redirect(w, r, r.URL.Path, http.StatusNotFound)
		err := global.Templates.ExecuteTemplate(w, "404", infoMsg)
		if err != nil {
			fmt.Println("deletePoll: problem parsing 404 template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	// if logged in user is not author of the post
	if user.Username != poll.Author {
		fmt.Println("Currently logged in user is not poll author")
		err := global.Templates.ExecuteTemplate(w, "403", infoMsg)
		if err != nil {
			fmt.Println("deletePoll: problem executing template: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	// everything is allright  delete poll with pollid
	_, err = global.DB.Exec(`DELETE from poll where id = $1`, poll.ID)
	if err != nil {
		fmt.Println("deletePoll:", poll.ID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// inform user about successful poll delete
	err = global.Templates.ExecuteTemplate(w, "info", infoMsg)
	if err != nil {
		fmt.Println("pollDelete: executing info template err:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
