package binance

type Interval string

const (
	InstitutionIdentifier = "BINANCE"
	maxRetries            = 3
	maxLimit              = 500

	baseUrl = "https://api.binance.com"
	// Endpoints
	getCandleStick  = "/api/v3/klines"
	getExchangeInfo = "/api/v3/exchangeInfo"
)
