package binance

type Interval string

const (
	BINANCE    = "BINANCE"
	maxRetries = 3
	maxLimit   = 500
	// calls per second
	rateLimit = 10

	baseUrl = "https://api.binance.com"
	// Endpoints
	getCandleStick  = "/api/v3/klines"
	getExchangeInfo = "/api/v3/exchangeInfo"
)
