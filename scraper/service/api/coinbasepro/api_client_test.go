package coinbasepro

import (
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoinbaseProClient(t *testing.T) {
	// TODO: Use DI instead of calling GetSecrets directly
	config.LoadEnv()
	secrets, _ := config.GetSecrets()
	exchangeClient := NewCoinbaseProAPIClient(secrets)
	pass := true
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := maxLimit * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength - time.Minute)
		candleStickData, err := exchangeClient.getCandleStickData(
			60,
			startTime,
			endTime,
			"BTC-USD",
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
	}) && pass
	pass = t.Run("TestGetProducts", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getProducts()
		assert.NoError(t, err)
		assert.NotNil(t, exchangeInfo)
		//fmt.Print(utils.PrettyJSON(exchangeInfo))
	}) && pass
	// Interface Methods
	pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
		pairs, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
		// TODO: This test is not scale-able unless we use DI to inject what pairs are supported
		expectedPairs := map[string]models.Symbol{
			"BTC-USD": {
				RawBase:         "BTC",
				RawQuote:        "USD",
				NormalizedBase:  "BTC",
				NormalizedQuote: "USD",
				ProductID:       "BTC-USD",
			},
			"ETH-USD": {
				RawBase:         "ETH",
				RawQuote:        "USD",
				NormalizedBase:  "ETH",
				NormalizedQuote: "USD",
				ProductID:       "ETH-USD",
			},
			"ETH-BTC": {
				RawBase:         "ETH",
				RawQuote:        "BTC",
				NormalizedBase:  "ETH",
				NormalizedQuote: "BTC",
				ProductID:       "ETH-BTC",
			},
		}
		assert.Equal(t, 3, len(pairs))
		for _, pair := range pairs {
			expectedSymbol, ok := expectedPairs[pair.ProductID]
			assert.Equal(t, true, ok)
			assert.Equal(t, expectedSymbol, *pair)
		}
	}) && pass

	// Should get all prices from [start, end)
	pass = t.Run("TestGetAllOHLCMarketData", func(t *testing.T) {
		expectedLength := 12000 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength)

		candleStickData, err := exchangeClient.GetAllOHLCMarketData(
			models.Symbol{
				ProductID: "BTC-USD",
			},
			time.Minute,
			startTime,
			endTime,
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, startTime.String(), candleStickData[0].StartTime.UTC().String())
		assert.Equal(t, endTime.String(), candleStickData[len(candleStickData)-1].EndTime.UTC().String())
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
	}) && pass

	assert.Equal(t, true, pass)
}
