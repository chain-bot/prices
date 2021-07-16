package coinbasepro

type Interval string

const (
	COINBASEPRO = "COINBASEPRO"
	maxLimit    = 300
	// calls per second

	baseURL = "https://api.pro.coinbase.com"
	// Endpoints
	getExchangeProducts = "/products"
	getCandles          = "/products/%s/candles"
)

var coinbaseProInstrumentFilter = map[string]bool{
	"USDT": true, // messes up tests and queries, most data is empty anyways
}
