package psql

import (
	"testing"

	"github.com/chain-bot/prices/app/configs"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	configs.LoadEnv()
	pass := true
	secrets, err := configs.GetSecrets()
	assert.NoError(t, err)
	var db *sqlx.DB

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
	}) && pass
	_, _ = db.Query("DROP DATABASE IF EXISTS ?", secrets.DBName)
	assert.Equal(t, true, pass)
}
