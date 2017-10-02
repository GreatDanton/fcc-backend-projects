package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql" // sql drivers

	"github.com/go-zoo/bone"
	_ "github.com/lib/pq" // importing postgres db drivers

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/controllers"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// main function for handling web application
func main() {
	config := global.ReadConfig()
	fmt.Printf("Starting server: http://127.0.0.1:%v\n", config.Port)
	r := bone.New()

	r.NotFoundFunc(controllers.Handle404)
	r.Get("/", http.HandlerFunc(controllers.FrontPage))
	r.HandleFunc("/register", controllers.Register)
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/poll/:poolID", controllers.ViewPoll)
	r.HandleFunc("/poll/:poolID/edit", controllers.EditPollHandler)
	r.HandleFunc("/new", controllers.CreateNewPoll)
	r.HandleFunc("/u/:userID", controllers.UserDetails)
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
