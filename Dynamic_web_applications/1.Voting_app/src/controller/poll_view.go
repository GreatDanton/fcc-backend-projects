package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

func ViewPoll(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayPoll(w, r, model.Poll{})
	case "POST":
		r.ParseForm()
		method := r.Form["_method"][0]
		switch method {
		case "post":
			postVote(w, r) // user posted vote
		case "delete":
			deletePoll(w, r) // user wanted to delete poll
		default:
			postVote(w, r)
		}
	default:
		displayPoll(w, r, model.Poll{})
	}
}

// displayPoll is handling GET request for VIEW POLL function
// displayPoll displays data for chosen poll /poll/:id and returns
//404 page if poll does not exist
func displayPoll(w http.ResponseWriter, r *http.Request, pollMsg model.Poll) {
	pollID := utilities.GetURLSuffix(r)

	poll, err := model.GetPollDetails(pollID)
	if err != nil {
		fmt.Println("Error while getting poll details: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	votes, err := model.GetPollVotes(pollID)
	if err != nil {
		fmt.Printf("Error while getting poll votes count: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	poll.Votes = votes
	user := model.LoggedInUser(r)
	poll.LoggedInUser = user

	// check if poll title exists and display relevant template with
	// poll data filled in
	if len(poll.Title) > 0 && len(poll.Options) > 0 {
		// TODO: fix this ugly implementation
		poll.Errors.PostVoteError = pollMsg.Errors.PostVoteError
		err = templates.Execute(w, "details", poll)
		if err != nil {
			fmt.Println("displayPoll: template err:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else { // if db does not return any rows -> poll does not exist, display 404
		err := templates.Execute(w, "404", poll)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}
