package coinbasepro

type Interval string

const (
	COINBASEPRO = "COINBASEPRO"
	maxRetries  = 3
	maxLimit    = 300
	// calls per second
	rateLimit = 2

	baseURL = "https://api.pro.coinbase.com"
	// Endpoints
	getExchangeProducts = "/products"
	getCandles          = "/products/%s/candles"
)
