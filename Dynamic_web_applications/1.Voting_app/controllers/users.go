package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// UserDetails is displaying details of chosen user
// details are: username and created pools
func UserDetails(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		userDetailsGET(w, r)
	case "POST":
		fmt.Println("Posting stuff in /u/")
	default:
		userDetailsGET(w, r)
	}
}

// User is used to display user details in /u/username
type User struct {
	Username string
}

// userDetailsGet renders userDetail template and displays users data
// username and created pools
func userDetailsGET(w http.ResponseWriter, r *http.Request) {
	user := User{}
	user.Username = strings.Split(r.URL.EscapedPath(), "/")[2]

	t := template.Must(template.ParseFiles("templates/users.html",
		"templates/navbar.html", "templates/styles.html"))
	err := t.Execute(w, user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
