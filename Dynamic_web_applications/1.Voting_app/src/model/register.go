package model

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// RegisterErrors stores error messages that will be displayed
// on register template if they occur
type RegisterErrors struct {
	Username      string
	ErrorUsername string

	Email      string
	ErrorEmail string

	Password      string
	ErrorPassword string
	LoggedInUser  User
}

// RegisterNewUser creates new row in database for new user
func RegisterNewUser(username string, passwordHash []byte, email string) (string, error) {
	// register our new user
	var id string
	err := global.DB.QueryRow(`INSERT into users(username, password_hash, email)
									values($1, $2, $3) RETURNING id`, username, passwordHash, email).Scan(&id)
	if err != nil {
		fmt.Printf("registerNewUser: problem inserting new user: %v\n", err)
		return "", fmt.Errorf("registerNewUser: problem inserting new user: %v", err)
	}
	return id, nil
}

// UserExistCheck checks if user with username exists
// false => user does not exist
// true => user already in database
// error => an error happened
func UserExistCheck(username string) (bool, error) {
	var usr string
	err := global.DB.QueryRow(`SELECT username from users
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

// UserEmailCheck checks if email already exists in database
// false => email does not exist
// true => email does exist in db
func UserEmailCheck(email string) (bool, error) {
	var e string
	err := global.DB.QueryRow(`SELECT email from users
							   where email = $1`, email).Scan(&e)
	if err != nil {
		// email does not exist
		if err == sql.ErrNoRows {
			return false, nil
		}
		// an error occurs on db lookup
		return true, err
	}
	// email does exist in database
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
