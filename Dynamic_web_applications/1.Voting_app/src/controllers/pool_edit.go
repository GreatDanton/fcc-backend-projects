package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

//TODO: FINISH HERE
// handling edit pool get requests
func EditPoolHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		editPoolView(w, r)
	case "POST":
		editPool(w, r)
	default:
		editPoolView(w, r)
	}
}

// editPool handles edit button press on poolDetails page
func editPoolView(w http.ResponseWriter, r *http.Request) {
	// insert stuff into fields
	//
	loggedUser := LoggedIn(r)
	if !loggedUser.LoggedIn {
		err := global.Templates.ExecuteTemplate(w, "403", http.StatusForbidden)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	poolID := strings.Split(r.URL.Path, "/")[2] //ustrings.Split(u, "/")[2]
	fmt.Println(poolID)
	pool, err := getPoolDetails(poolID)
	if err != nil {
		fmt.Println(err)
		return
	}
	pool.LoggedInUser = loggedUser
	if loggedUser.Username != pool.Author {
		fmt.Println("Currently logged in user is not the author")
		err := global.Templates.ExecuteTemplate(w, "403", pool)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	err = global.Templates.ExecuteTemplate(w, "editPool", pool)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

// TODO: finish this
//handling post request of editPoolHandler
func editPool(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)

	vars := mux.Vars(r)
	fmt.Println(vars)
	fmt.Println("POSTED STUFF")
}

// parsePoolParams fetches data from editPool/newPool form and returns:
// poolTitle, [voteOptions], error
func parsePoolParams(w http.ResponseWriter, r *http.Request) (string, []string, error) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", []string{}, err
	}

	poolTitle := strings.TrimSpace(r.Form["poolTitle"][0])
	// check if poolTitle exists else return template with error message
	if len(poolTitle) < 1 {
		e := newPoolError{TitleError: "Please add title to your pool"}
		err := global.Templates.ExecuteTemplate(w, "newPool", e)
		if err != nil {
			fmt.Println(err)
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
		voteOption := strings.TrimSpace(option[0])     // trim empty space from pool option
		if key != "poolTitle" && len(voteOption) > 0 { // filter out empty fields and title
			order = append(order, key)
		}
	}
	// if there are not at least 2 options to vote for return error into template
	if len(order) < 2 {
		e := newPoolError{Title: poolTitle, VoteOptionsError: "Please add at least two options"}
		err := global.Templates.ExecuteTemplate(w, "newPool", e)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return "", []string{}, err
		}
		return "", []string{}, fmt.Errorf("User added less than 2 vote options")
	}

	// this ensures pool options are inserted into database in
	// the same order as the end-user intended
	sort.Strings(order)
	voteOptions := make([]string, 0, len(order))
	for _, value := range order {
		voteOptions = append(voteOptions, value)
	}

	return poolTitle, voteOptions, nil
}
