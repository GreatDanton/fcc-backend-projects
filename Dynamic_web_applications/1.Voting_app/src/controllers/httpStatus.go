package controllers

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

type navbar404 struct {
	LoggedInUser User
}

// Handle404 handles httpNotFound response for bone router
func Handle404(w http.ResponseWriter, r *http.Request) {
	user := LoggedIn(r)
	navbar := navbar404{LoggedInUser: user}
	err := global.Templates.ExecuteTemplate(w, "404", navbar)
	if err != nil {
		fmt.Println("Handle404: problem parsing template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
