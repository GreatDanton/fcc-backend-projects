package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

// Register is handling registration of Voting application
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		registerDisplay(w, r, model.RegisterErrors{})
	case "POST":
		registerNewUser(w, r)
	default:
		registerDisplay(w, r, model.RegisterErrors{})
	}
}

// registerDisplay displays register template and possible error messages to the
// end user
func registerDisplay(w http.ResponseWriter, r *http.Request, errMsg model.RegisterErrors) {
	user := model.LoggedInUser(r)
	if user.LoggedIn {
		fmt.Println("registerDisplay: user is already logged in")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.Execute(w, "register", errMsg)
	if err != nil {
		fmt.Printf("registerDisplay: %v \n", err)
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

	errMsg := model.RegisterErrors{
		Username: username,
		Email:    email,
		Password: password,
	}

	// if passwords do not match, inform user and rerender template
	if password != passConfirm {
		errMsg.ErrorPassword = "Passwords do not match"
		registerDisplay(w, r, errMsg)
		return
	}

	// check if username already exist
	exist, err := model.UserExistCheck(username)
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
		registerDisplay(w, r, errMsg)
		return
	}

	// check if email already exists in database, emails should be unique
	//TODO: send confirmation email after registration
	exist, err = model.UserEmailCheck(email)
	if err != nil {
		fmt.Println("Register: userEmailCheck:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// email exists in database display error message
	if exist {
		fmt.Println("Register: userEmailCheck: Email is already registered")
		errMsg.ErrorEmail = "Email is already registered"
		registerDisplay(w, r, errMsg)
		return
	}

	// hash user inserted password
	passwordHash, err := model.HashPassword(password)
	if err != nil {
		fmt.Printf("HashPassword: %v \n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, err := model.RegisterNewUser(username, passwordHash, email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// create user session, log user in
	err = model.CreateUserSession(w, userID, username)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("/u/%v", username)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
