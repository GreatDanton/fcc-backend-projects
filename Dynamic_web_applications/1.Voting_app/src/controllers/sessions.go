package controllers

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// User struct used for creating tokens and parsing data from cookies
type User struct {
	ID       string
	Username string
	LoggedIn bool
}

// CreateToken creates signed token string with user id and username as payload
// signed string password is parsed from config.json file
func CreateToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = user.ID //User{ID: user.ID, Username: user.Username}
	claims["username"] = user.Username
	claims["loggedIn"] = user.LoggedIn
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims
	tokenEncode := []byte(global.Config.JWTtokenPassword)

	signedString, err := token.SignedString(tokenEncode)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

// CreateCookie creates cookie out of inserted arguments (user.id, user.username)
// and returns error when that is not possible
func CreateCookie(id, username string) (http.Cookie, error) {
	// set cookie and redirect
	expiration := time.Now().Add(7 * 24 * time.Hour) // cookie expires in 1 week
	u := User{ID: id, Username: username, LoggedIn: true}
	tokenString, err := CreateToken(u)
	if err != nil {
		e := fmt.Errorf("CreateCookie: CreateToken:%v", err)
		return http.Cookie{}, e
	}
	cookie := http.Cookie{Name: "GoVote", Value: tokenString,
		Expires: expiration, Path: "/", HttpOnly: true}

	return cookie, err
}

// DestroyCookie since we can't delete cookie on all browsers,
// it sets value of authentication cookie to blank and add expiration date = now
func DestroyCookie(r *http.Request) (http.Cookie, error) {
	_, err := r.Cookie("GoVote")
	// cookie does not exist
	if err != nil {
		return http.Cookie{}, err
	}
	c := http.Cookie{Name: "GoVote", Value: "", Expires: time.Now(), Path: "/", HttpOnly: true}
	return c, nil
}

// LoggedIn checks if cookie is present in client request and returns:
// true: if user is logged in
// false: if user is not logged in (cookie is not present, or an error occured)
func LoggedIn(r *http.Request) User {
	cookie, err := r.Cookie("GoVote")
	if err != nil {
		return User{}
	}
	tokenString := cookie.Value
	u, loggedIn := GetUserData(tokenString)
	if !loggedIn {
		return User{}
	}
	return u
}

// GetUserData gets userData from JWT token string and returns
// UserData (users.id, users.username)
// loggedIn (bool): true if token is formed correctly
//				  false if token is forged or an error occured
func GetUserData(tokenString string) (User, bool) {
	tokenEncode := []byte(global.Config.JWTtokenPassword)

	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return tokenEncode, nil
	})
	u := User{}

	// if error occured
	if err != nil {
		return u, false
	}

	// check if token is valid
	if !token.Valid {
		return u, false
	}

	id := token.Claims.(jwt.MapClaims)["id"]
	username := token.Claims.(jwt.MapClaims)["username"]
	loggedIn := token.Claims.(jwt.MapClaims)["loggedIn"]

	// type assertion, checking if returned values are actually strings
	// if they are not return empty user struct
	if idStr, ok := id.(string); ok {
		u.ID = idStr
	} else {
		return u, false
	}
	if usernameStr, ok := username.(string); ok {
		u.Username = usernameStr
	} else {
		return u, false
	}
	if loggedInStr, ok := loggedIn.(bool); ok {
		u.LoggedIn = loggedInStr
	} else {
		return u, false
	}

	return u, true
}
