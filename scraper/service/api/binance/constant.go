package binance

const (
	BINANCE    = "BINANCE"
	maxRetries = 3
	maxLimit   = 1000
	// calls per second
	rateLimit = 10

	baseUrl = "https://api.binance.com"
	// Endpoints
	getCandleStick  = "/api/v3/klines"
	getExchangeInfo = "/api/v3/exchangeInfo"

	klineSocketStream = "wss://stream.binance.com:9443/ws/%s@kline_%s"
)
