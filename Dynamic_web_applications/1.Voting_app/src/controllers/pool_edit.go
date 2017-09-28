package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// editPool handles edit button press on poolDetails page
func editPool(w http.ResponseWriter, r *http.Request) {
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
