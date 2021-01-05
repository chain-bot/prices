package config

import "os"

type Keys struct {
	BinanceApiKey string
	KrakenApiKey  string
	KucoinApiKey  string
}

func GetKeys() Keys {
	return Keys{
		BinanceApiKey: os.Getenv("BINANCE_API_KEY"),
	}
}
