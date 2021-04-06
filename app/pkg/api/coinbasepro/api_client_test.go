package coinbasepro

import (
	"github.com/mochahub/coinprice-scraper/app/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoinbaseProClient(t *testing.T) {
	exchangeClient := NewCoinbaseProAPIClient()
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
		expectedPair := models.Symbol{
			RawBase:         "BTC",
			RawQuote:        "USD",
			NormalizedBase:  "BTC",
			NormalizedQuote: "USD",
			ProductID:       "BTC-USD",
		}
		assert.GreaterOrEqual(t, len(pairs), 3)
		assert.Contains(t, pairs, &expectedPair)
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
