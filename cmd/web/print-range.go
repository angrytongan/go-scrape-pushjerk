package main

import (
	"fmt"
	"html/template"
	"net/http"
	"pj/internal/metricise"
)

const (
	sqlWorkouts = `
SELECT
	id,
	content
FROM
	posts
WHERE
	date >= ?
	AND date <= ?
ORDER BY
	date ASC
`
)

func (app *Application) PrintRange(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	pageData := map[string]any{}

	if len(values) == 0 {
		app.render(w, "print-range", pageData, http.StatusOK)

		return
	}

	errs := map[string]string{}
	workouts := []struct {
		ID      string
		Content template.HTML
	}{}
	fields := map[string]string{}

	start := values.Get("start")
	finish := values.Get("finish")

	if start == "" || finish == "" {
		if start == "" {
			errs["Start"] = "Need start date."
		}

		if finish == "" {
			errs["Finish"] = "Need finish date."
		}
	} else {
		rows, err := app.db.Query(sqlWorkouts, start, finish)
		if err != nil {
			errs["Query"] = fmt.Sprintf("app.db.Query(%s, %s): %v", start, finish, err)
		} else {
			for rows.Next() {
				var workout struct {
					ID      string
					Content template.HTML
				}

				if err := rows.Scan(&workout.ID, &workout.Content); err != nil {
					errs["Results"] = fmt.Sprintf("rows.Scan(): %v", err)
					break
				} else {
					workout.Content = template.HTML(metricise.Metricise(string(workout.Content)))
					workouts = append(workouts, workout)
				}
			}
		}
	}

	fields["Start"] = start
	fields["Finish"] = finish

	pageData["Errors"] = errs
	pageData["Workouts"] = workouts
	pageData["Fields"] = fields

	app.render(w, "print-range", pageData, http.StatusOK)
}
