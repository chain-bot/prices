package psql

import (
	"fmt"
	"github.com/chain-bot/prices/app/configs"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(
	secrets *configs.Secrets,
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
