package ftx

type Interval string

const (
	FTX      = "FTX"
	maxLimit = 5000
	// calls per second
	callsPerSecond = 30

	// Endpoints
	baseUrl             = "https://ftx.com/api"
	getHistoricalPrices = "/markets/%s/candles?resolution=%d&limit=%d&start_time=%d&end_time=%d"
	getMarkets          = "/markets"
)

var ftxToCoinprice = map[string]string{
	"USD": "USDT",
}

var coinpriceToftx = map[string]string{
	"USDT": "USD",
}
