package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

// User is used to display user details in /u/username
type userDetails struct {
	Username     string
	Polls        []model.Poll
	LoggedInUser model.User
	/* 	Pagination   pagination */
}

// UserDetails renders userDetail template and displays users data
// username and created polls
func UserDetails(w http.ResponseWriter, r *http.Request) {
	user := userDetails{}
	user.Username = utilities.GetURLSuffix(r)
	user.LoggedInUser = model.LoggedInUser(r)

	exist, err := model.UserExistCheck(user.Username)
	if err != nil { // user does not exist
		fmt.Println("userExistCheck:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// if user does not exist, display 404 page
	if !exist {
		fmt.Println("User does not exist")
		err = templates.Execute(w, "404", nil)
		if err != nil {
			fmt.Println("UserDetails:", err)
			return
		}
		return
	}
	limit := 20
	maxID := 0
	// get polls from user
	userPolls, err := model.GetUserPolls(user.Username, maxID, limit)
	if err != nil {
		fmt.Printf("getUserPoll: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user.Polls = userPolls
	/* 	p := handlePollPagination(r, maxID, userPolls, limit)
	   	user.Pagination = p */

	err = templates.Execute(w, "users", user)
	if err != nil {
		fmt.Println("UserDetails:", err)
		return
	}
}
