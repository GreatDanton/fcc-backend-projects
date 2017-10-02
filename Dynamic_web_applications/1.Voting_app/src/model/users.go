package model

import (
	"database/sql"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// User struct used for creating tokens and parsing data from cookies
type User struct {
	ID       string // user id
	Username string // users username
	LoggedIn bool   // true => logged in, false => not logged in
}

// LoggedInUser checks if GoVote(auth) cookie is present in client request and
// returns userStruct with {ID: user.id, Username: user.username, LoggedIn: true/false}
func LoggedInUser(r *http.Request) User {
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

// GetUserPolls fetches the database and returns user polls based on the
// the function arguments
// userName: user useranem
// maxID: maximum poll id
// limit: number of polls
func GetUserPolls(username string, maxID int, limit int) ([]Poll, error) {
	polls := []Poll{}
	var (
		id     string
		title  string
		author string
		time   time.Time
		votes  string
	)

	var rows *sql.Rows
	var err error
	if maxID == 0 {
		rows, err = global.DB.Query(`SELECT poll.id, poll.title, poll.created_by, poll.time, count(vote.poll_id)
							   		 from poll
							   		 LEFT JOIN users
							   		 on users.id = poll.created_by
							   		 LEFT JOIN vote
							   		 on vote.poll_id = poll.id
							   		 WHERE users.username = $1
							   		 GROUP BY poll.id
							   		 order by poll.id desc
							   		 limit $2`, username, limit)
	} else {
		// poll id can't be < 0
		rows, err = global.DB.Query(`SELECT poll.id, poll.title, poll.created_by, poll.time, count(vote.poll_id)
									 from poll
									 LEFT JOIN users
									 on users.id = poll.created_by
									 LEFT JOIN vote
									 on vote.poll_id = poll.id
									 WHERE users.username = $1
									 AND poll.id <= $2
									 GROUP BY poll.id
									 order by poll.id desc
									 limit $3`, username, maxID, limit)

	}

	if err != nil {
		return polls, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &time, &votes)
		if err != nil {
			return polls, err
		}
		timeAgo := utilities.TimeDiff(time) // create submitted ...ago string
		polls = append(polls, Poll{ID: id, Title: title, Author: author, Time: timeAgo, NumOfVotes: votes})
	}
	err = rows.Err()
	if err != nil {
		return polls, err
	}

	return polls, nil
}
