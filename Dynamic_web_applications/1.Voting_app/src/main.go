package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql" // sql drivers

	_ "github.com/lib/pq" // importing postgres db drivers

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/controllers"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/src/global"
)

// Config struct for holding data from config.json file

// main function for handling web application
func main() {
	config := global.ReadConfig()
	fmt.Printf("Starting server: http://127.0.0.1:%v\n", config.Port)

	http.HandleFunc("/", controllers.FrontPage)
	http.HandleFunc("/register/", controllers.Register)
	http.HandleFunc("/login/", controllers.Login)
	http.HandleFunc("/logout/", controllers.Logout)
	http.HandleFunc("/pool/", controllers.ViewPool)
	http.HandleFunc("/new/", controllers.CreateNewPool)
	http.HandleFunc("/u/", controllers.UserDetails)
	// handle public files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

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
