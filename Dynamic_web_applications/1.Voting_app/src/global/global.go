// Package global package is used for storing global variables that will be used
// across Voting application
package global

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
)

// Templates variable is storing templates on application startup
var Templates *template.Template

//DB global variable for storing database connection
var DB *sql.DB

// Config - global variable for providing global configuration
var Config Configuration

// Configuration is used store values from parsed config.json file
type Configuration struct {
	Port             string
	DbUser           string
	DbPassword       string
	DbName           string
	JWTtokenPassword string
}

// ReadConfig reads configuration file and exits if it does not exist or
// is wrongly formatter
func ReadConfig() Configuration {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Please add config.json file:", err)
	}
	config := Configuration{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal("Please format configuration file correctly:", err)
	}

	// Parses templates on config read and add them to global Template variable
	ParseTemplates()

	Config = config // add parsed json.config to global Config variable
	return config
}

// ParseTemplates parses templates folder and assigns parsed template to
// the global.Template variable
func ParseTemplates() {
	parsedTemplates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedTemplates, err)
	Templates = tmpl // assign parsed templates into global variable
}
