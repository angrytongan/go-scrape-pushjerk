package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"pj/pkg/routes"

	_ "github.com/mattn/go-sqlite3"
)

const port = 4000

func main() {
	// Database setup.
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println("couldn't open database", err)
		return
	}
	defer db.Close()

	// Router setup.
	r := routes.New(db)

	// Run application.
	fmt.Printf("listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
