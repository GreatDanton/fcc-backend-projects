package controller

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/model"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/templates"
)

// Login is handling login part of the application, by displaying
// login screen, validating inputs and taking care of user sessions
func Login(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		displayLogin(w, r, model.LoginErrors{})
	case "POST":
		login(w, r)
	default:
		displayLogin(w, r, model.LoginErrors{})
	}

}

// displayLogin displays login template and possible error messages / labels
func displayLogin(w http.ResponseWriter, r *http.Request, errMsg model.LoginErrors) {
	user := model.LoggedInUser(r)
	if user.LoggedIn {
		fmt.Println("displayLogin: User already logged in as", user.Username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	errMsg.LoggedInUser = user
	err := templates.Execute(w, "login", errMsg)
	if err != nil {
		fmt.Println("displayLogin:", err)
		return
	}
}

// logIn handles user login and displaying error messages
func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	username := strings.TrimSpace(r.Form["username"][0])
	password := r.Form["password"][0]
	errMsg := model.LoginErrors{Username: username}

	LoginUser, err := model.GetUserLoginData(username, password)
	if err != nil {
		if err == model.ErrUserDoesNotExist {
			errMsg.ErrorUsername = fmt.Sprintf("%v", model.ErrUserDoesNotExist)
			displayLogin(w, r, errMsg)
			return
		}
		// an actual error occured
		fmt.Println("login:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// compare password hashes
	bytePass := []byte(password)
	err = bcrypt.CompareHashAndPassword(LoginUser.PasswordHash, bytePass)
	if err != nil {
		fmt.Println("login: Wrong password")
		errMsg.ErrorPassword = "Wrong password"
		displayLogin(w, r, errMsg)
		return
	}

	// if password is correc create user session and login user
	err = model.CreateUserSession(w, LoginUser.ID, LoginUser.Username)
	if err != nil {
		fmt.Println("Login: CreateUserSession:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// User is logged in => redirect to front page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout is handling user logout by destroying user session
// or redirecting to login page if client is not logged in
func Logout(w http.ResponseWriter, r *http.Request) {
	user := model.LoggedInUser(r)
	if !user.LoggedIn { // if user is already logged in
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := model.DestroyUserSession(w, r)
	if err != nil {
		fmt.Println("Logout: DestroyUserSession:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, "logout", nil)
	if err != nil {
		fmt.Println("Logout:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
