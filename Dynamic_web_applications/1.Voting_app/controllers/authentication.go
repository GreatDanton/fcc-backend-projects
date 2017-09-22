package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/global"
)

// Register is handling registration of Voting application
func Register(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		registerGET(w, r, registerErrors{})
	case "POST":
		registerNewUser(w, r)
	default:
		registerGET(w, r, registerErrors{})
	}
}

type registerErrors struct {
	Username      string
	ErrorUsername string

	Email      string
	ErrorEmail string

	Password      string
	ErrorPassword string
}

// registerGET displays register template
func registerGET(w http.ResponseWriter, r *http.Request, errors registerErrors) {
	t := template.Must(template.ParseFiles("templates/register.html", "templates/navbar.html", "templates/styles.html"))
	err := t.Execute(w, errors)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form["email"][0]
	pass := r.Form["password"][0]
	passConfirm := r.Form["password-confirm"][0]
	username := r.Form["register-name"][0]

	// check password
	// check email
	//
	// check if username or email exist in db

	if pass != passConfirm {
		errMsg := registerErrors{
			Username:      username,
			Email:         email,
			Password:      pass,
			ErrorPassword: "Passwords do not match",
		}
		registerGET(w, r, errMsg)
		return
	}

	fmt.Println(username, email, pass, passConfirm)
}

// Login is handling login part + authentication
// TODO:
func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/login.html",
			"templates/navbar.html", "templates/styles.html"))
		err := t.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		// 1. Get data from form
		// 2. Check if user exists
		// 3. Encrypt password and compare it to hash
		// 4. Login + sessions
		// 5. ??
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		username := r.Form["username"][0]
		password := r.Form["password"][0]

		var (
			user     string
			passHash string
		)

		err = global.DB.QueryRow(`SELECT username, password_hash from users
										WHERE username = $1`, username).Scan(&user, &passHash)
		if err != nil {
			if err == sql.ErrNoRows { // if no rows exist
				fmt.Println("This user does not exist")
				return
			}
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if username != user {
			// inform user and return
			return
		}

		// hash password
		hashedPassword := password
		if hashedPassword != password {
			// inform user and return
			return
		}

		// handle session
		fmt.Fprintf(w, fmt.Sprintf("Logged in as %v", username))

	}
}
