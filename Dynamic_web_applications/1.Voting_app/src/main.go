package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql" // sql drivers

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // importing postgres db drivers

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/controllers"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// main function for handling web application
func main() {
	config := global.ReadConfig()
	fmt.Printf("Starting server: http://127.0.0.1:%v\n", config.Port)
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.FrontPage)
	r.HandleFunc("/register/", controllers.Register)
	r.HandleFunc("/login/", controllers.Login)
	r.HandleFunc("/logout/", controllers.Logout)
	r.HandleFunc("/poll/{pollID}", controllers.ViewPoll)
	r.HandleFunc("/poll/{pollID}/edit", controllers.EditPollHandler)
	r.HandleFunc("/new/", controllers.CreateNewPoll)
	r.HandleFunc("/u/{userID}", controllers.UserDetails)
	// handle public files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	http.Handle("/", r)

	// open connection with database using the fields from config
	connection := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	global.DB = db // assign to global db variable

	// setUp our database -> remove old tables and setup new ones
	//model.SetUpDB()

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
