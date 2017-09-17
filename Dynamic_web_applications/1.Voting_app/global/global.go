// Package global package is used for storing global variables that will be used
// across Voting application
package global

import (
	"database/sql"
)

//DB global variable for storing database connection
var DB *sql.DB
