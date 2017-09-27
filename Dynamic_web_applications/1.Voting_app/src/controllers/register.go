package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
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
	LoggedInUser  User
}

// registerGET displays register template and possible error messages to the user
func registerGET(w http.ResponseWriter, r *http.Request, errMsg registerErrors) {
	user := LoggedIn(r)
	errMsg.LoggedInUser = user

	err := global.Templates.ExecuteTemplate(w, "register", errMsg)
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
	exist, err := userExistCheck(username)
	// actual database error occured
	if err != nil {
		fmt.Printf("userExistCheck: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// user already exists
	if exist { // exist == true
		errMsg.ErrorUsername = "Username already taken"
		fmt.Println("Username already exists")
		registerGET(w, r, errMsg)
		return
	}

	// hash user inserted password
	passwordHash, err := HashPassword(password)
	if err != nil {
		fmt.Printf("HashPassword: %v \n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// register our new user
	var id string
	err = global.DB.QueryRow(`INSERT into users(username, password_hash, email)
								values($1, $2, $3) RETURNING id`, username, passwordHash, email).Scan(&id)
	if err != nil {
		fmt.Printf("registerNewUser: problem inserting new user: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// set up session
	cookie, err := CreateCookie(id, username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)
	url := fmt.Sprintf("/u/%v", username)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// check if user exists and return error if it does or an error happened
// false => user does not exist
// true => user already in database
func userExistCheck(username string) (bool, error) {
	var usr string
	err := global.DB.QueryRow(`SELECT users.username from users
							   where username = $1`, username).Scan(&usr)
	if err != nil {
		// if user does not exist, db returns 0 rows
		// we register him in outer function
		if err == sql.ErrNoRows {
			return false, nil
		}
		// if an actual error happens on db lookup, return err
		return true, err
	}
	// user exists
	return true, nil
}

// HashPassword hashes inserted users password
func HashPassword(password string) ([]byte, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}