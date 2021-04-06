package configs

// TODO: Move these to a table in PSQL
func GetSupportedAssets() map[string]bool {
	return map[string]bool{
		"BTC":  true,
		"ETH":  true,
		"USDT": true,
		"USD":  true,
	}
}
