package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase(t *testing.T) {
	pass := true
	pass = t.Run("TestNewDatabase", func(t *testing.T) {
		_, err := NewDatabase()
		assert.NoError(t, err)
	}) && pass
	assert.Equal(t, true, pass)
}
