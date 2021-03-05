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

// https://www.okex.com/docs/en/#spot-line_history
func getSupportedMap() map[string]bool {
	return map[string]bool{
		"BTC-USDT": true,
		"ETH-USDT": true,
		"LTC-USDT": true,
		"ETC-USDT": true,
		"XRP-USDT": true,
		"EOS-USDT": true,
		"BCH-USDT": true,
		"BSV-USDT": true,
		"TRX-USDT": true,
	}
}
