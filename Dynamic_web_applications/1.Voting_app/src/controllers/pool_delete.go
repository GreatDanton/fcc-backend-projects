package controllers

import (
	"fmt"
	"net/http"
)

// deletePool handles deleting chosen pools
func deletePool(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete pool")
}
