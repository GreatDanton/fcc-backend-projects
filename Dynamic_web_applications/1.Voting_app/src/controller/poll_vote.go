package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

// postVote function handles posting votes on each /poll/:id
func postVote(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in, if it's not return 403 forbidden
	user := model.LoggedInUser(r)
	if !user.LoggedIn {
		http.Redirect(w, r, r.URL.Path, http.StatusForbidden)
		return
	}
	r.ParseForm()
	pollID := utilities.GetURLSuffix(r)
	//pollID := strings.Split(r.URL.EscapedPath(), "/")[2]
	// get optionID, if the user did not pick anything
	// optionID is empty string
	var optionID string
	for key, value := range r.Form {
		if key == "voteOption" {
			optionID = value[0]
		}
	}

	pollMsg := model.Poll{}
	// if no vote option was chosen rerender template and display
	// error message to user
	if optionID == "" {
		pollMsg.Errors.PostVoteError = "Please pick your vote option"
		fmt.Println("postVote: no vote option was chosen")
		displayPoll(w, r, pollMsg)
		return
	}

	// check if user is changing vote options via html, this prevents
	// spamming votes for options that do not exist for this pollID
	voteOptions, err := model.GetVoteOptions(pollID)
	if err != nil {
		fmt.Println("postVote:", "getVoteOptions:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ok := utilities.StringInSlice(optionID, voteOptions)
	if !ok {
		pollMsg.Errors.PostVoteError = "You'll have to be more clever."
		fmt.Println("PostVote:", "User is changing vote options")
		displayPoll(w, r, pollMsg)
		return
	}
	// use user id of logged in user
	userID := user.ID

	err = model.PollAddUserVote(pollID, optionID, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// upon successful post request refresh page
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
