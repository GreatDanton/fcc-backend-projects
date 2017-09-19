package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"database/sql" // sql drivers

	_ "github.com/lib/pq" // importing postgres db drivers

	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/controllers"
	"github.com/greatdanton/fcc-backend-projects/Dynamic_web_applications/1.Voting_app/global"
)

// Config struct for holding data from config.json file
type Config struct {
	Port       string
	DbUser     string
	DbPassword string
	DbName     string
}

// readConfig reads configuration file and exits if it does not exist or
// is wrongly formatter
func readConfig() Config {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println(err)
		log.Fatal("Please add config.json file")
	}
	config := Config{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Println(err)
		log.Fatal("Please format configuration file correctly")
	}
	return config
}

// main function for handling web application
func main() {
	config := readConfig()
	fmt.Printf("Starting server: http://127.0.0.1:%v\n", config.Port)

	http.HandleFunc("/", controllers.FrontPage)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/view/", controllers.ViewPool)
	// handle public files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// open connection with database using the fields from config
	connection := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	global.DB = db // assign to global variable

	// setUp our database -> remove old tables and setup new ones
	//model.SetUpDB()

	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
