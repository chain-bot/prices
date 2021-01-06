package config

import "os"

type APIKeys struct {
	BinanceApiKey string
	KrakenApiKey  string
	KucoinApiKey  string
}
type DatabaseCredentials struct {
	User     string
	Password string
	DBName   string
	Port     int
	Host     string
}

type Secrets struct {
	APIKeys
	DatabaseCredentials
}

func GetSecrets() *Secrets {
	return &Secrets{
		APIKeys: APIKeys{
			BinanceApiKey: os.Getenv("BINANCE_API_KEY"),
		},
		DatabaseCredentials: DatabaseCredentials{
			User:     os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DATABASE"),
			Port:     5432,
			Host:     os.Getenv("POSTGRESQL_HOST"),
		},
	}
}
