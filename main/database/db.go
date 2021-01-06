package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/main/config"
)

type Database struct {
	*sqlx.DB
}

func NewDatabase() (*Database, error) {
	secrets := config.GetSecrets()
	psqlInfo := ""
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		secrets.DatabaseCredentials.Host,
		secrets.DatabaseCredentials.Port,
		secrets.DatabaseCredentials.User,
		secrets.DatabaseCredentials.Password,
		secrets.DatabaseCredentials.DBName)
	db, err := sqlx.Connect(databaseDriver, psqlInfo)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}
