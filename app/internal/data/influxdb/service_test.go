package influxdb

import (
	"github.com/chain-bot/prices/app/configs"
	"github.com/docker/distribution/context"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfluxDB(t *testing.T) {
	configs.LoadEnv()
	pass := true
	pass = t.Run("TestNewInfluxDBClient", func(t *testing.T) {
		// TODO: Use Uber fx
		secrets, _ := configs.GetSecrets()
		client, err := NewInfluxDBClient(secrets)
		assert.NoError(t, err)
		ctx := context.Background()
		healthCheck, err := (*client).Health(ctx)
		assert.NoError(t, err)
		assert.Equal(t, domain.HealthCheckStatusPass, healthCheck.Status)
	}) && pass
	assert.Equal(t, true, pass)
}
