package database

import (
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase(t *testing.T) {
	pass := true
	pass = t.Run("TestNewDatabase", func(t *testing.T) {
		// TODO: Use Uber fx
		secrets, _ := config.GetSecrets()
		_, err := NewDatabase(secrets)
		assert.NoError(t, err)
	}) && pass
	assert.Equal(t, true, pass)
}
