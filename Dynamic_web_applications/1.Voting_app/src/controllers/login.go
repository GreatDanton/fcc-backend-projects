package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// Login is handling login part of the application, by displaying
// login screen, validating inputs and taking care of user sessions
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

// used for displaying errors on login screen as well
// as rendering different view for loggedInUsers
type loginErrors struct {
	Username      string
	ErrorUsername string
	ErrorPassword string
	LoggedInUser  User
}

// displayLogIn function is used for displaying log in screen and
//handling error messages on login attempt
func displayLogIn(w http.ResponseWriter, r *http.Request, errMsg loginErrors) {
	user := LoggedIn(r)
	// if user is already logged in, redirect to front page
	if user.LoggedIn {
		fmt.Println("displayLogin: user is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	errMsg.LoggedInUser = user
	err := global.Templates.ExecuteTemplate(w, "login", errMsg)
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

	username := strings.TrimSpace(r.Form["username"][0])
	password := r.Form["password"][0]

	errMsg := loginErrors{Username: username}
	var (
		id       string
		user     string
		passHash []byte
	)

	err = global.DB.QueryRow(`SELECT id, username, password_hash from users
							  WHERE username = $1`, username).Scan(&id, &user, &passHash)

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
	// everything is allright, log in user
	err = createUserSession(id, username, w)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// User is logged in => redirect to front page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout destroys GoVote cookie and with that the user is logged out.
func Logout(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in, if not don't handle his request
	user := LoggedIn(r)
	if !user.LoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// user is logged in, destroy userSession and log him out
	err := destroyUserSession(w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// redirect to thank you for using our product page
	err = global.Templates.ExecuteTemplate(w, "logout", nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
