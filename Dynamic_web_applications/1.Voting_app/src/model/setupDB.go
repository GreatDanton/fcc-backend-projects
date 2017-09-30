package model

import (
	"fmt"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/controllers"

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// SetUpDB creates new database and fill it with some data
func SetUpDB() {
	// creating new tables in exact order
	createUserTable()
	createPollTable()
	createPollOption()
	createVoteTable()
}

// creates new user table and fill it with two users
func createUserTable() {
	_, err := global.DB.Exec("drop table if exists users cascade")
	if err != nil {
		fmt.Printf("Dropping users error: %v \n", err)
		return
	}

	_, err = global.DB.Exec(`Create table users(id serial primary key,
													username varchar(50) unique,
													email text unique,
													password_hash varchar(72) NOT NULL
													)`)
	if err != nil {
		fmt.Println("Creating Users Table Error:")
		fmt.Println(err)
	}

	passHash, err := controllers.HashPassword("bla bla")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = global.DB.Exec(`insert into users(username, email, password_hash)
										values('User1', 'email1@mail.com', $1)`, passHash)
	if err != nil {
		fmt.Println("Creating first user error:")
		fmt.Println(err)
	}

	passHash, err = controllers.HashPassword("bla2")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = global.DB.Exec(`insert into users(username, email, password_hash)
										values('User2', 'email2@mail.com', $1)`, passHash)
	if err != nil {
		fmt.Println("Creating second user error:")
		fmt.Println(err)
	}
}

// creates poll table and fill it with data
func createPollTable() {
	_, err := global.DB.Exec("drop table if exists poll cascade")
	if err != nil {
		fmt.Printf("Error dropping poll table: %v\n", err)
		return
	}

	_, err = global.DB.Exec(`Create table poll(id serial primary key,
												created_by integer references users(id) on delete cascade,
												time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
												title text)`)
	if err != nil {
		fmt.Printf("Error while creating poll table: %v\n", err)
		return
	}

	_, err = global.DB.Exec(`insert into poll(created_by, title) values(1, 'First title')`)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = global.DB.Exec(`insert into poll(created_by, title) values(2, 'Second title')`)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// create pollOptions and fill it with data
func createPollOption() {
	_, err := global.DB.Exec("drop table if exists pollOption cascade")
	if err != nil {
		fmt.Printf("Error dropping pollOption table: %v", err)
		return
	}

	_, err = global.DB.Exec(`create table pollOption(id serial primary key,
													poll_id integer references Poll(id) on delete cascade,
													option text)`)
	if err != nil {
		fmt.Printf("Error while creating pollOption table: %v\n", err)
		return
	}

	_, err = global.DB.Exec(`insert into pollOption(poll_id, option) values (1, 'First option')`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = global.DB.Exec(`insert into pollOption(poll_id, option) values (1, 'Second option')`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = global.DB.Exec(`insert into pollOption(poll_id, option) values (2, 'third option')`)
	if err != nil {
		fmt.Println(err)
	}
	_, err = global.DB.Exec(`insert into pollOption(poll_id, option) values (2, 'fourth option')`)
	if err != nil {
		fmt.Println(err)
	}
}

// create vote table
func createVoteTable() {
	_, err := global.DB.Exec("drop table if exists vote cascade")
	if err != nil {
		fmt.Printf("Error while dropping vote table: %v\n", err)
		return
	}

	_, err = global.DB.Exec(`create table vote(id serial,
											   poll_id integer references Poll(id) on delete cascade,
											   option_id integer references pollOption(id) on delete cascade,
											   voted_by integer references users(id) on delete cascade);`)

	if err != nil {
		fmt.Printf("Error while creating vote table: %v\n", err)
	}

	_, err = global.DB.Exec(`insert into vote(poll_id, option_id, voted_by)
										values(1, 1, 1)`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = global.DB.Exec(`insert into vote(poll_id, option_id, voted_by)
										values(1, 1, 2)`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = global.DB.Exec(`insert into vote(poll_id, option_id, voted_by)
										values(2, 3, 1)`)
	if err != nil {
		fmt.Println(err)
	}
}
