package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/config"
)

func NewDatabase(
	secrets *config.Secrets,
) (*sqlx.DB, error) {
	psqlInfo := ""
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		secrets.DatabaseCredentials.Host,
		secrets.DatabaseCredentials.Port,
		secrets.DatabaseCredentials.User,
		secrets.DatabaseCredentials.Password,
		secrets.DatabaseCredentials.DBName)
	return sqlx.Connect(databaseDriver, psqlInfo)
}
