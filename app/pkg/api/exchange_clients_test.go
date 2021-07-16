package api

import (
	"context"
	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"testing"
	"time"
)

func InjectDependencies(
	test interface{},
) error {
	fxApp := fx.New(
		GetAPIProviders(),
		fx.Provide(configs.GetSecrets),
		fx.Invoke(
			test,
		),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		return err
	}
	return fxApp.Stop(context.Background())
}

func TestExchangeClients(t *testing.T) {
	configs.LoadEnv()
	err := InjectDependencies(func(clients ExchangeClients) {
		for i := range clients.Clients {
			exchangeClient := clients.Clients[i]
			t.Run("TestGetSupportedPairs", func(t *testing.T) {
				symbols, err := exchangeClient.GetSupportedPairs()
				assert.Nil(t, err, exchangeClient.GetExchangeIdentifier())
				assert.NotEmpty(t, symbols, exchangeClient.GetExchangeIdentifier())
			})
			t.Run("TestGetExchangeIdentifier", func(t *testing.T) {
				identifier := exchangeClient.GetExchangeIdentifier()
				assert.NotEmpty(t, identifier, exchangeClient.GetExchangeIdentifier())
			})
			t.Run("TestGetAllOHLCVMarketData", func(t *testing.T) {
				expectedLength := 12000 * time.Minute
				startTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
				endTime := startTime.Add(expectedLength)
				pairs, _ := exchangeClient.GetSupportedPairs()
				assert.NotEmpty(t, pairs, exchangeClient.GetExchangeIdentifier())
				pair := filterPairsForBTC(pairs)
				assert.NotEmpty(t, pair, exchangeClient.GetExchangeIdentifier(), pairs)
				candleStickData, err := exchangeClient.GetAllOHLCVMarketData(
					*pair,
					time.Minute,
					startTime,
					endTime,
				)
				assert.NoError(t, err, exchangeClient.GetExchangeIdentifier())
				assert.NotEmpty(t, candleStickData, exchangeClient.GetExchangeIdentifier())
				assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData), exchangeClient.GetExchangeIdentifier())
				assert.Equal(t, startTime.String(), candleStickData[len(candleStickData)+1].StartTime.UTC().String(), exchangeClient.GetExchangeIdentifier())
				assert.Equal(t, endTime.String(), candleStickData[len(candleStickData)-1].EndTime.UTC().String(), exchangeClient.GetExchangeIdentifier())
			})
		}
	})
	assert.NoError(t, err)
}

// BTC is the oldest market (has the most data); safest for testing
func filterPairsForBTC(pairs []*models.Symbol) *models.Symbol {
	for i := range pairs {
		if pairs[i].NormalizedBase == "BTC" {
			return pairs[i]
		}
	}
	return nil
}
