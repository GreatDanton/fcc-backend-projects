package controller

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

// CreateNewPoll takes care of handling creation of the new poll in url: /new
// add post title to database
func CreateNewPoll(w http.ResponseWriter, r *http.Request) {
	user := model.LoggedInUser(r)
	// check if user is logged in, otherwise redirect to /login page
	if !user.LoggedIn {
		http.Redirect(w, r, "/login/", http.StatusSeeOther)
		return
	}

	poll := model.Poll{LoggedInUser: user}
	//errMsg := newPollError{LoggedInUser: user}

	if r.Method == "GET" {
		err := global.Templates.ExecuteTemplate(w, "newPoll", poll)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		pollTitle, voteOptions, err := parsePollParams(w, r, "newPoll")
		if err != nil {
			// displaying error message is done in function
			fmt.Println("CreateNewPoll:", "parsePollParams:", err)
			return
		}

		pollOptions := make([]string, 0, len(voteOptions))
		for _, value := range voteOptions {
			option := r.Form[value][0] // text of the voteOption
			pollOptions = append(pollOptions, option)
		}

		pollID, err := model.AddNewPoll(pollTitle, pollOptions, user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		url := fmt.Sprintf("/poll/%v", pollID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

//
// parsePollParams fetches data from editPoll/newPoll templates form and returns:
// pollTitle, [voteOptions], error
func parsePollParams(w http.ResponseWriter, r *http.Request, template string) (string, []string, error) {
	if template == "edit" {
		template = "editPoll"
	} else {
		template = "newPoll"
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", []string{}, err
	}
	errMsg := model.Poll{}

	pollTitle := strings.TrimSpace(r.Form["pollTitle"][0])
	// check if pollTitle exists else return template with error message
	if len(pollTitle) < 1 {
		errMsg.Errors.TitleError = "Please add title to your poll"
		err := templates.Execute(w, template, errMsg)
		if err != nil {
			fmt.Println("parsePollParams:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return "", []string{}, err
		}
		return "", []string{}, fmt.Errorf("Title of the post is missing")
	}

	order := make([]string, 0, len(r.Form))
	// r.Form returns a map, we have to add fields in db in correct order
	//  (=> that is in the same order the user wanted to post options)
	// so we don't confuse the end user, why their options are borked
	for key, option := range r.Form {
		voteOption := strings.TrimSpace(option[0])     // trim empty space from poll option
		if key != "pollTitle" && len(voteOption) > 0 { // filter out empty fields and title
			order = append(order, key)
		}
	}
	// if there are not at least 2 options to vote for return error into template
	if len(order) < 2 {
		errMsg.Errors.Title = pollTitle
		errMsg.Errors.VoteOptionsError = "Please add at least two options"
		// add vote options to the poll struct, otherwise options are missing upon
		// template rerender
		errMsg.Errors.VoteOptions = order
		err := templates.Execute(w, template, errMsg)
		if err != nil {
			fmt.Println("parsePollParams:", err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return "", []string{}, err
		}
		return "", []string{}, fmt.Errorf("User added less than 2 vote options")
	}

	// this ensures poll options are inserted into database in
	// the same order as the end-user intended
	sort.Strings(order)
	voteOptions := make([]string, 0, len(order))
	for _, value := range order {
		voteOptions = append(voteOptions, value)
	}

	return pollTitle, voteOptions, nil
}
