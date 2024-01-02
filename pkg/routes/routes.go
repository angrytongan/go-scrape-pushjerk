package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type MyHandler struct {
	DB *sql.DB
}

func New(db *sql.DB) chi.Router {
	m := &MyHandler{db}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", m.Home)

	return r
}

func (m *MyHandler) Home(w http.ResponseWriter, r *http.Request) {
	row := m.DB.QueryRow("SELECT COUNT(id) FROM posts")

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
