package controllers

import (
	"fmt"
	"net/http"
)

// Register is handling registration of Voting application
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Register page")
}

// Login is handling login part
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login page")
}
