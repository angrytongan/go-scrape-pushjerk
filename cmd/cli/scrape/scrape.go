package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlDropTable = `DROP TABLE IF EXISTS posts`

	sqlCreateTable = `CREATE TABLE posts (id TEXT, title TEXT, date DATE, content TEXT)`

	sqlInsert = `INSERT INTO posts (id, title, date, content) VALUES (?, ?, ?, ?)`

	layoutDate = "Mon, Jan _2, 2006"
)

type Application struct {
	Writer func(*colly.HTMLElement)
	db     *sql.DB
}

func (app *Application) SQLWriter(e *colly.HTMLElement) {
	id := e.Attr("id")
	title := e.ChildText("h2.entry-title")

	date, err := time.Parse(layoutDate, title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "time.Parse(%s): failed conversion: %v\n", title, err)
		date = time.Now()
	}

	content, _ := e.DOM.Html()

	_, err = app.db.Exec(sqlInsert, id, title, date, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to insert record", err)
	}
}

func (app *Application) TSVWriter(e *colly.HTMLElement) {
	id := e.Attr("id")
	title := e.ChildText("h2.entry-title")

	date, err := time.Parse(layoutDate, title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "time.Parse(%s): failed conversion: %v\n", title, err)
		date = time.Now()
	}

	content := strings.ReplaceAll(e.ChildText("div.entry-content"), "\n", " ")

	fmt.Printf("%s\t%s\t%s\t%s\n", id, title, date, content)
}

func New(tsv, sqlite bool) *Application {
	app := &Application{}
	if sqlite {
		db, err := sql.Open("sqlite3", "database.db")
		if err != nil {
			panic(fmt.Sprintf("couldn't open database: %v", err))
		}

		_, err = db.Exec(sqlDropTable)
		if err != nil {
			panic(fmt.Sprintf("failed to drop database: %v", err))
		}
		_, err = db.Exec(sqlCreateTable)
		if err != nil {
			panic(fmt.Sprintf("failed to create table: %v", err))
		}

		app.db = db
		app.Writer = app.SQLWriter
	} else if tsv {
		fmt.Printf("ID\ttitle\tdate\tcontent\n")

		app.Writer = app.TSVWriter
	}

	return app
}

func (app *Application) Close() {
	if app.db != nil {
		_ = app.db.Close()
	}
}

func main() {
	tsv := flag.Bool("tsv", false, "dump the scrape as a tsv")
	sqlite := flag.Bool("sqlite", false, "dump to sqlite database")
	flag.Parse()

	if !*tsv && !*sqlite {
		fmt.Fprintf(os.Stderr, "require either tsv or sqlite flags\n")
		return
	}

	app := New(*tsv, *sqlite)

	// Setup colly.
	c := colly.NewCollector(
		colly.Async(true),
	)

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*pushjerk*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add limit to collector: %v\n", err)
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Fprintln(os.Stderr, "visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Fprintln(os.Stderr, "something went wrong:", err)
	})

	// Setup database.
	// Parse the pager.
	c.OnHTML("a.next", func(r *colly.HTMLElement) {
		l := r.Attr("href")
		if l == "" {
			return
		}

		tokens := strings.Split(strings.Trim(l, "/"), "/")
		if len(tokens) == 0 {
			fmt.Fprintln(os.Stderr, "no tokens to parse")
			return
		}

		next := tokens[len(tokens)-1]
		_, err := strconv.Atoi(next)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to get next page number", err)
			return
		}

		if err = c.Visit(l); err != nil {
			fmt.Fprintf(os.Stderr, "c.Visit(): %v\n", err)
		}
	})

	// Grab the workout and write it.
	c.OnHTML("article", app.Writer)

	if err = c.Visit("https://pushjerk.com"); err != nil {
		fmt.Fprintf(os.Stderr, "c.Visit(https://pushjerk.com): %v\n", err)
	}
	c.Wait()
}
