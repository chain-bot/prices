package configs

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

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

type ServerConfig struct {
	Port int
}

type Environment string

type Secrets struct {
	Environment
	ServerConfig
	DatabaseCredentials
	InfluxDbCredentials
}

func GetSecrets() (*Secrets, error) {
	LoadEnv()
	postgresPort, err := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
	if err != nil {
		log.WithFields(log.Fields{
			"err":             err.Error(),
			"POSTGRESQL_PORT": os.Getenv("POSTGRESQL_PORT"),
		}).Errorf("error getting POSTGRESQL_PORT, setting default value")
		postgresPort = 5432
	}
	influxDBPort, err := strconv.Atoi(os.Getenv("INFLUXDB_PORT"))
	if err != nil {
		log.WithFields(log.Fields{
			"err":           err.Error(),
			"INFLUXDB_PORT": os.Getenv("INFLUXDB_PORT"),
		}).Errorf("error getting INFLUXDB_PORT, setting default value")
		influxDBPort = 8086
	}
	serverPort, err := strconv.Atoi(os.Getenv("PRICES_API_PORT"))
	if err != nil {
		log.WithFields(log.Fields{
			"err":             err.Error(),
			"PRICES_API_PORT": os.Getenv("PRICES_API_PORT"),
		}).Errorf("error getting PRICES_API_PORT, setting default value")
		serverPort = 8080
	}
	return &Secrets{
		ServerConfig: ServerConfig{
			Port: serverPort,
		},
		Environment: Environment(os.Getenv("CHAINBOT_ENV")),
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

func (s *Secrets) IsLocal() bool {
	return s.Environment == LocalEnv || s.Environment == NilEnv
}
