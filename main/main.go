package main

import (
	"log"

	"github.com/mochahub/coinprice-scraper/main/app"
)

func main() {
	// TODO: Find a Better Logging Framework
	log.Println("Scraper Config:")
	app.StartScraper()
}
