package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
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
		fmt.Println("visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("something went wrong:", err)
	})

	// Setup database.
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println("couldn't open database", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS posts")
	if err != nil {
		fmt.Println("failed to drop database", err)
		return
	}
	_, err = db.Exec("CREATE TABLE posts (id TEXT, title TEXT, content TEXT)")
	if err != nil {
		fmt.Println("failed to create table")
		return
	}

	// Parse the pager.
	c.OnHTML("a.next", func(r *colly.HTMLElement) {
		l := r.Attr("href")
		if l == "" {
			return
		}

		tokens := strings.Split(strings.Trim(l, "/"), "/")
		if len(tokens) == 0 {
			fmt.Println("no tokens to parse")
			return
		}

		next := tokens[len(tokens)-1]
		_, err := strconv.Atoi(next)
		if err != nil {
			fmt.Println("failed to get next page number", err)
			return
		}
		c.Visit(l)
	})

	// Grab the workout and write to database.
	c.OnHTML("article", func(e *colly.HTMLElement) {
		id := e.Attr("id")
		title := e.ChildText("h2.entry-title")
		content, _ := e.DOM.Html()

		_, err := db.Exec("insert into posts (id, title, content) values (?, ?, ?)",
			id, title, content)
		if err != nil {
			fmt.Println("failed to insert record", err)
		}
	})

	c.Visit("https://pushjerk.com")
	c.Wait()
}
