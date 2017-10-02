package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// LoginErrors struct is used for displaying errors on login screen
type LoginErrors struct {
	Username      string
	ErrorUsername string
	ErrorPassword string
	LoggedInUser  User
}

// LoginUser is used to gather data from database query
type LoginUser struct {
	ID           string
	Username     string
	PasswordHash []byte
}

// ErrUserDoesNotExist error message if user does not exist in database
var ErrUserDoesNotExist = errors.New("This user does not exist")

// GetUserLoginData get's the login data for chosen username
func GetUserLoginData(username string, password string) (LoginUser, error) {
	var (
		id       string
		user     string
		passHash []byte
	)

	err := global.DB.QueryRow(`SELECT id, username, password_hash from users
							  WHERE username = $1`, username).Scan(&id, &user, &passHash)
	tryingUser := LoginUser{}

	if err != nil {
		if err == sql.ErrNoRows { // if no rows exist
			return tryingUser, ErrUserDoesNotExist
		}
		return tryingUser, fmt.Errorf("GetUserLoginData: A database error occured: %v", err)
	}
	tryingUser.ID = id
	tryingUser.Username = user
	tryingUser.PasswordHash = passHash

	return tryingUser, nil
}
