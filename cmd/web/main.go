package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const port = 4000

type Workout struct {
	ID      string
	Title   string
	Content string
}

type Application struct {
	tpl *template.Template
	db  *sql.DB
}

func New() *Application {
	tpl := template.Must(template.ParseGlob("templates/*.tmpl"))

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	return &Application{
		tpl: tpl,
		db:  db,
	}
}

func (app *Application) Close() {
	app.db.Close()
}

func (app *Application) render(w http.ResponseWriter, pageName string, pageData map[string]any, statusCode int) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := app.tpl.ExecuteTemplate(w, pageName, pageData); err != nil {
		panic(err)
	}
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.Query(`
	SELECT
		CAST(SUBSTR(id, 6) AS INTEGER) AS id,
		title
	FROM
		posts
	ORDER BY
		id DESC
	`)
	if err != nil {
		panic(err)
	}

	var workouts []Workout

	for rows.Next() {
		var workout Workout
		if err := rows.Scan(&workout.ID, &workout.Title); err != nil {
			panic(err)
		}
		workouts = append(workouts, workout)
	}

	pageData := map[string]any{
		"Workouts": workouts,
	}

	app.render(w, "home", pageData, http.StatusOK)
}

func (app *Application) Random(w http.ResponseWriter, r *http.Request) {
	row := app.db.QueryRow(`
		SELECT
			CAST(SUBSTR(id, 6) AS INTEGER) AS id
		FROM
			posts
		ORDER BY
			RANDOM()
		LIMIT
			1
	`)

	var id int
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/workout/"+strconv.Itoa(id), http.StatusFound)
}

func (app *Application) prePostWorkouts(id int) (int, int) {
	var preID, thisID, postID int
	q := fmt.Sprintf(`
		WITH
			parsed_posts AS (
				SELECT
					CAST(SUBSTR(posts.id, 6) AS INTEGER) AS id
				FROM
					posts
			),
			numbered_posts AS (
				SELECT
					id,
					row_number() OVER (ORDER BY id) as rownum
				FROM
					parsed_posts
			),
			current AS (
				SELECT
					rownum
				FROM
					numbered_posts
				WHERE
					id = %d
			)

		SELECT
			numbered_posts.id
		FROM
			numbered_posts, current
		WHERE
			ABS(numbered_posts.rownum - current.rownum) <= 1
		ORDER BY
			numbered_posts.rownum
	`, id)
	rows, err := app.db.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&preID); err != nil {
			panic(err)
		}
	}

	if rows.Next() {
		if err := rows.Scan(&thisID); err != nil {
			panic(err)
		}
	}

	if rows.Next() {
		if err := rows.Scan(&postID); err != nil {
			panic(err)
		}
	}

	if preID == id { // we have no previous record
		postID = thisID
		preID = 0
	}

	return preID, postID
}

func (app *Application) Workout(w http.ResponseWriter, r *http.Request) {
	var workout Workout

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}
	q := fmt.Sprintf("SELECT id, title, content FROM posts WHERE id = 'post-%d'", id)

	row := app.db.QueryRow(q) //, id)
	if err := row.Scan(&workout.ID, &workout.Title, &workout.Content); err != nil {
		panic(err)
	}

	preID, postID := app.prePostWorkouts(id)
	pageData := map[string]any{
		"Workout": workout,
		"PreID":   preID,
		"PostID":  postID,
	}
	app.render(w, "workout", pageData, http.StatusOK)
}

func main() {
	app := New()
	defer app.Close()

	// Router setup.
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", app.Home)
	r.Get("/random", app.Random)
	r.Get("/workout/{id}", app.Workout)

	// Run application.
	fmt.Printf("listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
