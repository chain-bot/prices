package psql

import (
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/app/configs"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDatabase(t *testing.T) {
	configs.LoadEnv()
	pass := true
	secrets, _ := configs.GetSecrets()
	var db *sqlx.DB
	var err error

	pass = t.Run("TestNewDatabase", func(t *testing.T) {
		// TODO: Use Uber fx
		db, err = NewDatabase(secrets)
		assert.NoError(t, err)
	}) && pass
	assert.Equal(t, true, pass)

	secrets.DBName = "test1"
	_, _ = db.Query("CREATE DATABASE ?", secrets.DBName)
	pass = t.Run("TestNewDatabase", func(t *testing.T) {
		version, err := RunMigrations(db, secrets)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, version)
		log.Printf("database version %d \n", version)
	}) && pass
	_, _ = db.Query("DROP DATABASE IF EXISTS ?", secrets.DBName)
	assert.Equal(t, true, pass)
}
