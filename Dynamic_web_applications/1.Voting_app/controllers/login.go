package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/global"
)

// login.go is used for displaying login screen, validating inputs
// and taking care of user sessions

// Login is handling login part + authentication
// TODO:
func Login(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		displayLogIn(w, r, loginErrors{})
	case "POST":
		logIn(w, r)
	default:
		displayLogIn(w, r, loginErrors{})
	}

}

type loginErrors struct {
	Username      string
	ErrorUsername string
	ErrorPassword string
}

// used for displaying log in screen and handling error messages
func displayLogIn(w http.ResponseWriter, r *http.Request, errMsg loginErrors) {
	t := template.Must(template.ParseFiles("templates/login.html",
		"templates/navbar.html", "templates/styles.html"))
	err := t.Execute(w, errMsg)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// logIn handles login protocol, sets up session and informs user about the
// wrong username&password combination
func logIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	username := r.Form["username"][0]
	password := r.Form["password"][0]

	errMsg := loginErrors{Username: username}
	var (
		user     string
		passHash []byte
	)

	err = global.DB.QueryRow(`SELECT username, password_hash from users
											WHERE username = $1`, username).Scan(&user, &passHash)

	if err != nil {
		if err == sql.ErrNoRows { // if no rows exist
			errMsg.ErrorUsername = "This user does not exist"
			displayLogIn(w, r, errMsg)
			fmt.Println("This user does not exist")
			return
		}
		fmt.Printf("logIn: db: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	bytePass := []byte(password)
	err = bcrypt.CompareHashAndPassword(passHash, bytePass)
	// if hashed db password and inserted password do not match,
	// display error message to the user
	if err != nil {
		fmt.Printf("Login compare hash passwords: %v\n", err)
		errMsg.ErrorPassword = "Password does not match"
		displayLogIn(w, r, errMsg)
		return
	}

	// password & username match
	// Login - TODO: create session
	fmt.Println("Logged in")
	url := fmt.Sprintf("/u/%v", username)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
