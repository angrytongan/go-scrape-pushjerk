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

type DBHandler struct {
	db *sql.DB
}

func New(db *sql.DB) chi.Router {
	h := &DBHandler{db}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", h.Home)
	r.Get("/random", h.Random)
	r.Get("/workout/{id}", h.Workout)

	return r
}

func (h *DBHandler) Home(w http.ResponseWriter, r *http.Request) {
	row := h.db.QueryRow("SELECT COUNT(id) FROM posts")

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

type Workout struct {
	id      string
	title   string
	content string
}

func (h *DBHandler) Random(w http.ResponseWriter, r *http.Request) {
	row := h.db.QueryRow("SELECT id, title, content FROM posts ORDER BY RANDOM() LIMIT 1")

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

func (h *DBHandler) Workout(w http.ResponseWriter, r *http.Request) {
	var workout Workout

	id := chi.URLParam(r, "id")
	q := fmt.Sprintf("SELECT id, title, content FROM posts WHERE id = 'post-%s'", id)

	fmt.Println(q)

	row := h.db.QueryRow(q, id)
	err := row.Scan(&workout.id, &workout.title, &workout.content)
	if err != nil {
		log.Printf("failed to scan workout: %w\n", err)
		http.Error(w, "failed to scan workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("%s\n%s\n%s", workout.id, workout.title, strings.ReplaceAll(workout.content, "\n", "<br />"))))
}
