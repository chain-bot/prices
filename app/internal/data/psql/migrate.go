package psql

import (
	"github.com/chain-bot/prices/app/configs"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func RunMigrations(
	db *sqlx.DB,
	secrets *configs.Secrets,
) (int, error) {
	if err := db.Ping(); err != nil {
		log.WithField("err", err.Error()).Fatalf("could not ping DB")
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{DatabaseName: secrets.DBName})
	if err != nil {
		log.WithField("err", err.Error()).Fatalf("could not get DB driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		configs.GetMigrationDir(), // file://path/to/directory
		secrets.DatabaseCredentials.DBName, driver)
	if err != nil {
		log.WithField("err", err.Error()).Fatalf("could not apply migrations")
		return 0, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.WithField("err", err.Error()).Fatalf("could not sync DB")
		return 0, err
	}
	version, _, err := m.Version()
	log.WithField("version", version).Info("database migrated")
	return int(version), err
}
