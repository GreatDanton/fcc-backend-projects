package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

// Handle404 handles default httpNotFound response for 'bone' router
func Handle404(w http.ResponseWriter, r *http.Request) {
	err := templates.Execute(w, "404", nil)
	if err != nil {
		fmt.Println("Handle404:", err)
		return
	}
}

// info struct for displaying correct navbar when executing
// info template
type info struct {
	LoggedInUser model.User
}
