package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type APIKeys struct {
	BinanceApiKey            string
	CoinbaseProApiKey        string
	CoinbaseProApiSecret     string
	CoinbaseProApiPassphrase string
	KrakenApiKey             string
	KucoinApiKey             string
	KucoinApiSecret          string
	KucoinApiPassphrase      string
}
type DatabaseCredentials struct {
	User     string
	Password string
	DBName   string
	Port     int
	Host     string
}

type InfluxDbCredentials struct {
	User     string
	Password string
	Token    string
	Port     int
	Host     string
	Org      string
	Bucket   string
}

type Secrets struct {
	APIKeys
	DatabaseCredentials
	InfluxDbCredentials
}

func GetSecrets() (*Secrets, error) {
	LoadEnv()
	postgresPort, err := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
	if err != nil {
		return nil, err
	}
	influxDBPort, err := strconv.Atoi(os.Getenv("INFLUXDB_PORT"))
	if err != nil {
		return nil, err
	}
	return &Secrets{
		APIKeys: APIKeys{
			BinanceApiKey:            os.Getenv("BINANCE_API_KEY"),
			CoinbaseProApiKey:        os.Getenv("COINBASE_PRO_API_KEY"),
			CoinbaseProApiSecret:     os.Getenv("COINBASE_PRO_API_SECRET"),
			CoinbaseProApiPassphrase: os.Getenv("COINBASE_PRO_API_KEY_PASSPHRASE"),
			KucoinApiKey:             os.Getenv("KUCOIN_API_KEY"),
			KucoinApiSecret:          os.Getenv("KUCOIN_API_SECRET"),
			KucoinApiPassphrase:      os.Getenv("KUCOIN_API_PASSPHRASE"),
		},
		DatabaseCredentials: DatabaseCredentials{
			User:     os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DATABASE"),
			Port:     postgresPort,
			Host:     os.Getenv("POSTGRESQL_HOST"),
		},
		InfluxDbCredentials: InfluxDbCredentials{
			User:     os.Getenv("INFLUXDB_ADMIN_USER"),
			Password: os.Getenv("INFLUXDB_ADMIN_USER_PASSWORD"),
			Token:    os.Getenv("INFLUXDB_ADMIN_USER_TOKEN"),
			Port:     influxDBPort,
			Host:     os.Getenv("INFLUXDB_HOST"),
			Org:      os.Getenv("INFLUXDB_ORG"),
			Bucket:   os.Getenv("INFLUXDB_BUCKET_CANDLE"),
		},
	}, nil
}
