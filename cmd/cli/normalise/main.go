package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName = "database.db"

	sqlAllWorkouts = `
		SELECT
			id, title, date, content
		FROM
			posts
		--WHERE
		  --ID IN ("post-8912")
		ORDER BY date DESC
	`
)

type Workout struct {
	ID      string
	Title   string
	Date    time.Time
	Content string
}

func run() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("sql.Open(%s): %w", dbName, err)
	}
	defer db.Close()

	rows, err := db.Query(sqlAllWorkouts)
	if err != nil {
		return fmt.Errorf("db.Query(): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var workout Workout
		if err := rows.Scan(
			&workout.ID,
			&workout.Title,
			&workout.Date,
			&workout.Content,
		); err != nil {
			return fmt.Errorf("rows.Scan(): %w", err)
		}

		cleaned := normalise(workout.Content)
		fmt.Printf("ID %-5s: ", workout.ID)
		for _, section := range cleaned {
			fmt.Print(section.Section, " ")
		}
		fmt.Println()
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
