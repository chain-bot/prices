package okex

type Interval string

const (
	OKEX     = "OKEX"
	maxLimit = 300
	// calls per second
	callsPerSecond = 10

	// Endpoints
	baseUrl        = "https://www.okex.com"
	getCandles     = "/api/spot/v3/instruments/%s/history/candles"
	getInstruments = "/api/spot/v3/instruments"
)
