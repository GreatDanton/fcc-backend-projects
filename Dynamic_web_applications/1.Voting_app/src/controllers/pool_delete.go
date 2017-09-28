package controllers

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// deletePool handles deleting chosen pools
func deletePool(w http.ResponseWriter, r *http.Request) {
	user := LoggedIn(r)
	// if user is not logged in render 403 template
	infoMsg := info{LoggedInUser: user}

	if !user.LoggedIn {
		fmt.Println("User is not logged in")
		err := global.Templates.ExecuteTemplate(w, "403", infoMsg)
		if err != nil {
			fmt.Println("delete Pool: Problem parsing templates:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	poolID := utilities.GetURLSuffix(r)
	pool, err := getPoolDetails(poolID)
	if err != nil {
		fmt.Println("Delete: getPoolDetails", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// if title is empty, pool does not exist
	if pool.Title == "" {
		// using post -> redirect, or post -> render?
		//http.Redirect(w, r, r.URL.Path, http.StatusNotFound)
		err := global.Templates.ExecuteTemplate(w, "404", infoMsg)
		if err != nil {
			fmt.Println("deletePool: problem parsing 404 template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	// if logged in user is not author of the post
	if user.Username != pool.Author {
		fmt.Println("Currently logged in user is not pool author")
		err := global.Templates.ExecuteTemplate(w, "403", infoMsg)
		if err != nil {
			fmt.Println("deletePool: problem executing template: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	// everything is allright  delete pool with poolid
	_, err = global.DB.Exec(`DELETE from pool where id = $1`, pool.ID)
	if err != nil {
		fmt.Println("deletePool:", pool.ID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// inform user about successful pool delete
	err = global.Templates.ExecuteTemplate(w, "info", infoMsg)
	if err != nil {
		fmt.Println("poolDelete: executing info template err:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
