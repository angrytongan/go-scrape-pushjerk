.PHONY: web scrape help

.DEFAULT_GOAL: help

help:
	@echo "make web - run web server"
	@echo "make scrape-tsv - run scraper and write to stdout as tsv"
	@echo "make scrape-db - run scraper and write to sqlite db"

web:
	# go run pj/cmd/web
	air

scrape-tsv:
	@go run cmd/cli/scrape.go --tsv

scrape-db:
	@go run cmd/cli/scrape.go --sqlite
