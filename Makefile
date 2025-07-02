.PHONY: web scrape help

.DEFAULT_GOAL: help

help: # me
	@grep -e "^[a-z0-9_-]*:.* # .*" Makefile | sed -e 's/\(^[a-z0-9_-]*\):.* # \(.*\)/\1: \2/'

web: # Run web application.
	# go run pj/cmd/web
	air

scrape-tsv: # Run scraper, writing to .tsv.
	@go run cmd/cli/scrape/scrape.go --tsv

scrape-db: # Run scraper, writing to sqlite db.
	@go run cmd/cli/scrape/scrape.go --sqlite

normalise: # Run normaliser.
	@go run cmd/cli/normalise/*
