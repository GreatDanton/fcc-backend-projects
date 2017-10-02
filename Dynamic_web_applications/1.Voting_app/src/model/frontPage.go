package model

import (
	"database/sql"
	"time"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"
)

// GetFrontPageData returns array of polls based on chosen maxID(max poll id) and limit of results
func GetFrontPageData(maxID int, limit int) ([]Poll, error) {
	polls := []Poll{}
	rows, err := fpQuery(maxID, limit)
	if err != nil {
		return polls, err
	}
	defer rows.Close()

	var (
		id         string
		title      string
		author     string
		time       time.Time
		numOfVotes string
	)

	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &time, &numOfVotes)
		if err != nil {
			return polls, err
		}
		// get time difference in human readable format
		t := utilities.TimeDiff(time)
		polls = append(polls, Poll{ID: id, Title: title,
			Author: author, Time: t, NumOfVotes: numOfVotes})
	}
	err = rows.Err()
	if err != nil {
		return polls, err
	}
	// if error does not happen, return results
	return polls, nil
}

// fpQuery picks the most suitable sql query based on maxID of poll and returns
// sql rows and error
func fpQuery(maxID int, limit int) (*sql.Rows, error) {
	if maxID == 0 { // maxID = 0 when we perform the first query on "/"
		rows, err := global.DB.Query(`SELECT poll.id, poll.title,
								  users.username, poll.time,
								  (select count(*) as votes from vote where vote.poll_id = poll.id)
								  FROM poll
								  LEFT JOIN users on users.id = poll.created_by
								  ORDER BY poll.id desc
								  limit $1`, limit)
		return rows, err
	}

	rows, err := global.DB.Query(`SELECT poll.id, poll.title,
								  users.username, poll.time,
								  (select count(*) as votes from vote where vote.poll_id = poll.id)
								  FROM poll
								  LEFT JOIN users on users.id = poll.created_by
								  WHERE poll.id <= $1
								  ORDER BY poll.id desc
								  limit $2`, maxID, limit)
	return rows, err
}
