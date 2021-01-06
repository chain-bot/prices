package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mochahub/coinprice-scraper/main/config"
	"log"
)

func RunMigrations(database *Database) error {
	secrets := config.GetSecrets()
	if err := database.DB.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}
	driver, err := postgres.WithInstance(database.DB.DB, &postgres.Config{DatabaseName: secrets.DBName})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}
	sourceURL := fmt.Sprintf("file://%s", migrationDir)
	println(sourceURL)
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL, // file://path/to/directory
		secrets.DatabaseCredentials.DBName, driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
		return err
	}
	log.Println("Database migrated")
	return nil
}
