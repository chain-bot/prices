package binance

const (
	InstitutionIdentifier = "BINANCE"
	// TODO(Zahin): We should switch between test and prod baseUrl in a config, not a const
	baseUrl = "https://api.binance.com"

	// Endpoints
	getCandleStick = "/api/v3/klines"

	maxLimit = 500
)
