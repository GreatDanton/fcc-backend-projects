package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

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

// registerError is used to populate user inserted fields and
// displaying errors to the user inside register template
type registerErrors struct {
	Username      string
	ErrorUsername string

	Email      string
	ErrorEmail string

	Password      string
	ErrorPassword string
}

// registerGET displays register template and possible error messages to the user
func registerGET(w http.ResponseWriter, r *http.Request, errMsg registerErrors) {
	t := template.Must(template.ParseFiles("templates/register.html", "templates/navbar.html", "templates/styles.html"))
	err := t.Execute(w, errMsg)
	if err != nil {
		fmt.Printf("registerGet: %v \n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// registerNewUser takes care of registering new users as well as
// backend user input validation
func registerNewUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := strings.TrimSpace(r.Form["register-name"][0])
	email := strings.TrimSpace(r.Form["email"][0])
	password := r.Form["password"][0]
	passConfirm := r.Form["password-confirm"][0]

	errMsg := registerErrors{
		Username: username,
		Email:    email,
		Password: password,
	}

	// TODO check if email exist in db
	//TODO: check email -> send confirmation email to registered address

	// if passwords do not match, inform user and rerender template
	if password != passConfirm {
		errMsg.ErrorPassword = "Passwords do not match"
		registerGET(w, r, errMsg)
		return
	}

	// check if username already exist
	err := userExistCheck(username)
	if err != nil {
		e := fmt.Sprint(err)
		// inform our user his username is already taken
		if e == "User already exists" {
			errMsg.ErrorUsername = "Username already taken"
			fmt.Println("Username already exists")
			registerGET(w, r, errMsg)
			return
		}
		// actual database error occured
		fmt.Printf("userExistCheck: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// hash user inserted password
	passwordHash, err := hashPassword(password)
	if err != nil {
		fmt.Printf("hashPassword: %v \n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// register our new user
	_, err = global.DB.Exec(`INSERT into users(username, password_hash, email)
							  values($1, $2, $3)`, username, passwordHash, email)
	if err != nil {
		fmt.Printf("registerNewUser: problem inserting new user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// if everything is allright, login user and redirect him to his page
	// login user: TODO
	url := fmt.Sprintf("/u/%v", username)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// used to parse user fields from db if userExists
type userCheck struct {
	Username string
	Email    string
}

// check if user exists and return error if it does or an error happened
func userExistCheck(username string) error {
	user := userCheck{}
	err := global.DB.QueryRow(`SELECT users.username from users
							   where username = $1`, username).Scan(&user.Username)
	if err != nil {
		// if user does not exist, db returns 0 rows
		// we register him in outer function
		if err == sql.ErrNoRows {
			return nil
		}
		// if an actual error happens on db lookup, return err
		return err
	}
	// user exists, return error
	return fmt.Errorf("%v", "User already exists")
}

// hashPassword hashes users password
func hashPassword(password string) ([]byte, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}
