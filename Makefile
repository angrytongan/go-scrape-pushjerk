.PHONY: web scrape help

help:
	@echo "make web - run web server"
	@echo "make scrape - run scraper"

web:
	# go run pj/cmd/web
	air

scrape:
	go run cmd/cli/scrape.go
