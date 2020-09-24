package main

import (
	"log"

	"github.com/mochahub/crypto-price-server/main/app"
)

func main() {
	// TODO: Find a Better Logging Framework
	log.Println("Scraper Config:")
	app.StartScraper()
}
