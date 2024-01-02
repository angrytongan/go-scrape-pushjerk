package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", Home(db))
	r.Get("/random", Random(db))

	return r
}

func Home(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRow("SELECT COUNT(id) FROM posts")

		var num int
		err := row.Scan(&num)
		if err != nil {
			log.Printf("failed to scan post count", err)
			http.Error(w, "failed to scan post count", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "text/html")

		w.Write([]byte(fmt.Sprintf("home page - %d workouts", num)))
	}
}

type Workout struct {
	id      string
	title   string
	content string
}

func Random(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRow("SELECT id, title, content FROM posts ORDER BY RANDOM() LIMIT 1")

		var workout Workout
		err := row.Scan(&workout.id, &workout.title, &workout.content)
		if err != nil {
			log.Printf("failed to scan workout", err)
			http.Error(w, "failed to scan workout", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("random - %s %s %s", workout.id, workout.title, strings.ReplaceAll(workout.content, "\n", "<br />"))))
	}
}
