package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // importing postgres db drivers

	"github.com/go-zoo/bone"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/controller"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// main function for handling web application
func main() {
	// reading configuration
	config := global.ReadConfig()
	fmt.Printf("Starting server: http://127.0.0.1:%v\n", config.Port)
	r := bone.New()

	/* 	r.NotFoundFunc(controllers.Handle404) */
	r.NotFoundFunc(controller.Handle404)
	r.Get("/", http.HandlerFunc(controller.FrontPage))
	r.HandleFunc("/register", controller.Register)
	r.HandleFunc("/new", controller.CreateNewPoll)
	r.HandleFunc("/u/:userID", controller.UserDetails)
	r.HandleFunc("/login", controller.Login)
	r.HandleFunc("/logout", controller.Logout)
	r.HandleFunc("/poll/:pollID", controller.ViewPoll)
	r.HandleFunc("/poll/:pollID/edit", controller.EditPollHandler)
	// handle public files
	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// open database connection using the fields from config
	connection := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	global.DB = db // assign to global db variable

	// setUp our database -> remove old tables and setup new ones
	//model.SetUpDB()

	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
