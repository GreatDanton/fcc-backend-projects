package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

// UserDetails is displaying details of chosen user
// details are: username and created pools
func UserDetails(w http.ResponseWriter, r *http.Request) {
	switch m := r.Method; m {
	case "GET":
		userDetailsGET(w, r)
	default:
		userDetailsGET(w, r)
	}
}

// User is used to display user details in /u/username
type userDetails struct {
	Username     string
	Pools        []userPool
	LoggedInUser User
}

// userDetailsGet renders userDetail template and displays users data
// username and created pools
func userDetailsGET(w http.ResponseWriter, r *http.Request) {
	user := userDetails{}
	user.Username = strings.Split(r.URL.EscapedPath(), "/")[2]
	u := LoggedIn(r)
	user.LoggedInUser = u

	userPools, err := getUserPools(user.Username)
	if err != nil {
		fmt.Printf("getUserPool: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Pools = userPools

	err = global.Templates.ExecuteTemplate(w, "users", user)
	if err != nil {
		fmt.Printf("userDetailsGET: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// userPools is used to display user pools in user details page
type userPool struct {
	ID         string
	Title      string
	Time       string
	NumOfVotes string
}

func getUserPools(username string) ([]userPool, error) {
	pool := []userPool{}
	var (
		id    string
		title string
		time  time.Time
		votes string
	)
	rows, err := global.DB.Query(`SELECT pool.id, pool.title, pool.time, count(vote.pool_id)
							   from pool
							   LEFT JOIN users
							   on users.id = pool.created_by
							   LEFT JOIN vote
							   on vote.pool_id = pool.id
							   WHERE users.username = $1
							   GROUP BY pool.id
							   order by pool.id desc
							   limit 20`, username)
	if err != nil {
		return pool, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &title, &time, &votes)
		if err != nil {
			return pool, err
		}
		timeAgo := utilities.TimeDiff(time) // create submitted ...ago string
		pool = append(pool, userPool{ID: id, Title: title, Time: timeAgo, NumOfVotes: votes})
	}
	err = rows.Err()
	if err != nil {
		return pool, err
	}

	return pool, err
}
