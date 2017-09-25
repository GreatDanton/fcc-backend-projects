package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
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
	LoggedInUser  User
}

// used for displaying log in screen and handling error messages
func displayLogIn(w http.ResponseWriter, r *http.Request, errMsg loginErrors) {
	user := LoggedIn(r)
	errMsg.LoggedInUser = user

	t := template.Must(template.ParseFiles("templates/login.html",
		"templates/navbar.html", "templates/styles.html"))
	err := t.ExecuteTemplate(w, "login", errMsg)
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

	// create cookie out id & username
	cookie, err := CreateCookie(id, username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie) // set cookie -> user is logged in

	url := fmt.Sprintf("/u/%v", username)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// Logout destroys GoVote cookie and with that the user is not logged in anymore
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := DestroyCookie(r)

	if err != nil {
		// no cookie is present but the user press logout - how to deal with this?
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// an actual error occured
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	// redirect to thank you for using our product site
	t := template.Must(template.ParseFiles("templates/logout.html", "templates/navbar.html", "templates/styles.html"))
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
