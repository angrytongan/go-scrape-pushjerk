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

type Application struct {
	Writer func(*colly.HTMLElement)
	db     *sql.DB
}

func (app *Application) SQLWriter(e *colly.HTMLElement) {
	id := e.Attr("id")
	title := e.ChildText("h2.entry-title")
	content, _ := e.DOM.Html()

	_, err := app.db.Exec("insert into posts (id, title, content) values (?, ?, ?)",
		id, title, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to insert record", err)
	}
}

func (app *Application) TSVWriter(e *colly.HTMLElement) {
	id := e.Attr("id")
	title := e.ChildText("h2.entry-title")
	content := strings.Replace(e.ChildText("div.entry-content"), "\n", " ", -1)

	fmt.Printf("%s\t%s\t%s\n", id, title, content)
}

func New(tsv, sqlite bool) *Application {
	app := &Application{}
	if sqlite {
		db, err := sql.Open("sqlite3", "database.db")
		if err != nil {
			panic(fmt.Sprintf("couldn't open database: %v", err))
		}

		_, err = db.Exec("DROP TABLE IF EXISTS posts")
		if err != nil {
			panic(fmt.Sprintf("failed to drop database: %v", err))
		}
		_, err = db.Exec("CREATE TABLE posts (id TEXT, title TEXT, content TEXT)")
		if err != nil {
			panic(fmt.Sprintf("failed to create table: %v", err))
		}

		app.db = db
		app.Writer = app.SQLWriter
	} else if tsv {
		fmt.Printf("ID\ttitle\tcontent\n")

		app.Writer = app.TSVWriter
	}

	return app
}

func (app *Application) Close() {
	if app.db != nil {
		app.db.Close()
	}
}

func main() {
	tsv := flag.Bool("tsv", false, "dump the scrape as a tsv")
	sqlite := flag.Bool("sqlite", false, "dump to sqlite database")
	flag.Parse()

	if *tsv == false && *sqlite == false {
		fmt.Fprintf(os.Stderr, "require either tsv or sqlite flags\n")
		return
	}

	app := New(*tsv, *sqlite)

	// Setup colly.
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*pushjerk*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

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
		c.Visit(l)
	})

	// Grab the workout and write it.
	c.OnHTML("article", app.Writer)

	c.Visit("https://pushjerk.com")
	c.Wait()
}
