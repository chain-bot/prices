package binance

const (
	BINANCE        = "BINANCE"
	maxLimit       = 1000
	callsPerMinute = 1200

	baseUrl = "https://api.binance.com"
	// Endpoints
	getCandleStick  = "/api/v3/klines"
	getExchangeInfo = "/api/v3/exchangeInfo"
)
