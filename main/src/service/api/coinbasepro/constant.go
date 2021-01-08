package coinbasepro

type Interval string

const (
	COINBASE   = "COINBASE"
	maxRetries = 3
	maxLimit   = 300
	// calls per second
	rateLimit = 2

	baseURL = "https://api.pro.coinbase.com"
	// Endpoints
	getExchangeProducts = "/products"
	getCandles          = "/products/%s/candles"
)

var coinbaseToCoinprice = map[string]string{
	"USD": "USDT",
}

var coinpriceToCoinbase = map[string]string{
	"USDT": "USD",
}
