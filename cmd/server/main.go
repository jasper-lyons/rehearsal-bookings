package main

import (
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	// GO Embed doesn't support embedding files from outside the module boundary
	// (in this case, anything outside of this directory) but we want to store
	// templates and static files at the route directory so we need to treat them
	// as their own go modules (you'll see the main.go files in the web/static and
	// web/templates directories) so that we can import them here and access the 
	// embedded static files!
	"rehearsal-bookings/web/static"
	"rehearsal-bookings/web/templates"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	// init the db
	db, err := sql.Open("sqlite3", "development.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
