package kucoin

type Interval string

const (
	KUCOIN   = "KUCOIN"
	maxLimit = 1500
	// calls per second
	callsPerMinute = 1800

	// Endpoints
	baseUrl    = "https://api.kucoin.com"
	getKlines  = "/api/v1/market/candles"
	getSymbols = "/api/v1/symbols"
)
