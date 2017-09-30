package controllers

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/utilities"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// Poll structure used to parse values from database
type Poll struct {
	ID      string
	Author  string     // author username
	Title   string     // title of the poll
	Options [][]string // contains [[option title, option_id]]
	Votes   [][]string // contains [[vote Option, vote count]]
	//ErrorPostVote string     // display error when user submits his vote
	LoggedInUser User       // User struct for rendering different templates based on user login status
	Errors       pollErrors // displaying error messages in new/edit poll templates
}

type pollErrors struct {
	Title            string // upon error fill input with .Title
	TitleError       string // display error when title is not suitable
	VoteOptions      []string
	VoteOptionsError string // display error when user submitted < 2 vote options
	PostVoteError    string // display error on user vote
	EditPollError    string // display error upon submitting edit poll form
}

// ViewPoll takes care for handling existing polls in /poll/poll_id
// displaying existing polls and handling voting part of the poll
func ViewPoll(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayPoll(w, r, Poll{})

	case "POST":
		r.ParseForm()
		method := r.Form["_method"][0]
		switch method {
		case "post":
			postVote(w, r) // user posted vote
		case "delete":
			deletePoll(w, r) // user wanted to delete poll
		default:
			postVote(w, r)
		}
	default:
		displayPoll(w, r, Poll{})
	}
}

// displayPoll is handling GET request for VIEW POOL function
// displayPoll displays data for chosen poll /poll/:id and returns
//404 page if poll does not exist
func displayPoll(w http.ResponseWriter, r *http.Request, pollMsg Poll) {
	pollID := utilities.GetURLSuffix(r)

	poll, err := getPollDetails(pollID)
	if err != nil {
		fmt.Println("Error while getting poll details: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	votes, err := getPollVotes(pollID)
	if err != nil {
		fmt.Printf("Error while getting poll votes count: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	poll.Votes = votes
	user := LoggedIn(r)
	poll.LoggedInUser = user

	// check if poll title exists and display relevant template with
	// poll data filled in
	if len(poll.Title) > 0 && len(poll.Options) > 0 {
		// TODO: fix this ugly implementation
		poll.Errors.PostVoteError = pollMsg.Errors.PostVoteError
		err = global.Templates.ExecuteTemplate(w, "details", poll)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else { // if db does not return any rows -> poll does not exist, display 404
		err := global.Templates.ExecuteTemplate(w, "404", poll)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// getPollDetails returns {pollID, Title, Author, [polloption, polloptionID]}
// from database for chosen pollID
func getPollDetails(pollID string) (Poll, error) {
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

// getPollVotes returns vote count for poll with chosen pollID
// returns [[Vote option 1, count 1], [Vote option 2, count 2]]
func getPollVotes(pollID string) ([][]string, error) {
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
