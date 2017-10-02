package model

import (
	"fmt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// Poll structure used to parse values from database
type Poll struct {
	ID         string
	Author     string // author username
	Title      string // title of the poll
	Time       string
	Options    [][]string // contains [[option title, option_id]]
	Votes      [][]string // contains [[vote Option, vote count]]
	NumOfVotes string     // for displaying number of all votes for the Poll
	//ErrorPostVote string     // display error when user submits his vote
	LoggedInUser User       // User struct for rendering different templates based on user login status
	Errors       PollErrors // displaying error messages in new/edit poll templates
}

// PollErrors are used to display /new/edit poll tempaltes
type PollErrors struct {
	Title            string // upon error fill input with .Title
	TitleError       string // display error when title is not suitable
	VoteOptions      []string
	VoteOptionsError string // display error when user submitted < 2 vote options
	PostVoteError    string // display error when user submits his vote
	EditPollError    string // display error upon submitting edit poll form
}

// GetPollDetails get's data about poll with id = pollID
// returned: {pollID, Title, Author, [pollOption, polloptionID]}
func GetPollDetails(pollID string) (Poll, error) {
	poll := Poll{}
	rows, err := global.DB.Query(`SELECT poll.id, poll.title, users.username, polloption.option,
								  polloption.id from poll
								  LEFT JOIN pollOption
								  on poll.id = pollOption.poll_id
								  LEFT JOIN users
								  on users.id = poll.created_by
								  where poll.id = $1;`, pollID)
	if err != nil {
		return poll, err
	}
	defer rows.Close()
	// defining variables for parsing rows from db
	var (
		id           string
		title        string
		author       string
		pollOption   string
		pollOptionID string
	)
	// parsing rows from database
	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &pollOption, &pollOptionID)
		if err != nil {
			return poll, err
		}
		poll.ID = id
		poll.Title = title
		poll.Author = author

		option := []string{pollOption, pollOptionID}
		poll.Options = append(poll.Options, option)
	}
	err = rows.Err()
	if err != nil {
		return poll, err
	}
	return poll, nil
}

// GetPollVotes returns vote count for poll with id = pollID
// returns [[VoteOption1, count1], [VoteOption2, count2]]
func GetPollVotes(pollID string) ([][]string, error) {
	votes := [][]string{} //Votes{}
	// returns: optionID, optionName, number of votes => sorted by increasing id
	// this ensures vote options results are returned the same way as they were posted
	rows, err := global.DB.Query(`SELECT pollOption.id, pollOption.option,
								  count(vote.option_id) from polloption
								  LEFT JOIN vote
								  on polloption.id = vote.option_id
								  where polloption.poll_id = $1
								  group by pollOption.id
								  order by pollOption.id asc`, pollID)
	if err != nil {
		return votes, err
	}
	defer rows.Close()

	var (
		id         string
		voteOption string
		count      string
	)
	for rows.Next() {
		err := rows.Scan(&id, &voteOption, &count)
		if err != nil {
			return votes, err
		}
		// appending results of table rows to Votes
		vote := []string{voteOption, count}
		votes = append(votes, vote)
	}
	err = rows.Err()
	if err != nil {
		return votes, err
	}
	return votes, nil
}

// DeletePoll deletes poll with id = pollID and
// returns error if error occur
func DeletePoll(pollID string) error {
	_, err := global.DB.Exec(`DELETE from poll where id = $1`, pollID)
	if err != nil {
		fmt.Println("DeletePoll: ID:", pollID, err)
		return err
	}
	return nil
}
