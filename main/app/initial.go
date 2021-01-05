package app

func InitScrapper() {
	// First Time Run Can Potentially take Hours (back filling market data)
	StartScraper()
	// Subsequent scheduled crons should take seconds
	StartScrapperCron()
}
