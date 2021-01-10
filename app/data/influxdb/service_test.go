package influxdb

import (
	"github.com/docker/distribution/context"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/mochahub/coinprice-scraper/app/utils"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfluxDB(t *testing.T) {
	utils.LoadEnv()
	pass := true
	pass = t.Run("TestNewInfluxDBClient", func(t *testing.T) {
		// TODO: Use Uber fx
		secrets, _ := config.GetSecrets()
		client, err := NewInfluxDBClient(secrets)
		assert.NoError(t, err)
		ctx := context.Background()
		healthCheck, err := (*client).Health(ctx)
		assert.NoError(t, err)
		assert.Equal(t, domain.HealthCheckStatusPass, healthCheck.Status)
	}) && pass
	assert.Equal(t, true, pass)
}
