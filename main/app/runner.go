package app

import (
	"log"
)

func StartScraper() {
	log.Println("Starting Scraper")
	defer log.Println("Stopping Scraper")

	// TODO: Move Below Comments to a Space Document
	// General Flow:
	// - Loop Through Implemented Exchanges
	// - Each Exchange is a goroutine
	// - Each goroutine writes to influx independently (conflicts will take latest price?)
	//	- Each exchange writes to their own bucket with expiry
	//	- A separate bucket aggregates price data from exchange buckets periodically
	// - TODO: Interface File that Each Exchange Inherits
	//	- Implement Influx Config and Functions here to abstract from exchange
}
