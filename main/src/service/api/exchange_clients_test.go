package api

import (
	"github.com/mochahub/coinprice-scraper/main/src/service/api/binance"
	"github.com/mochahub/coinprice-scraper/main/src/service/api/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestBinanceClient(t *testing.T) {
	clients := []ExchangeAPIClient{
		binance.NewBinanceAPIClient(os.Getenv("BINANCE_API_KEY")),
	}
	for _, exchangeClient := range clients {
		pass := true
		pass = t.Run("TestGetAllOHLCMarketData", func(t *testing.T) {
			expectedLength := 1000 * time.Minute
			startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			endTime := startTime.Add(expectedLength - time.Minute)
			candleStickData, err := exchangeClient.GetAllOHLCMarketData(
				"BTC",
				"USDT",
				common.Minute,
				startTime,
				endTime,
			)
			assert.NoError(t, err)
			assert.NotEmpty(t, candleStickData)
			assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
		}) && pass
		pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
			symbols, err := exchangeClient.GetSupportedPairs()
			assert.Nil(t, err)
			assert.NotEmpty(t, symbols)
			//fmt.Print(utils.PrettyJSON(symbols))
		}) && pass
		pass = t.Run("TestGetExchangeIdentifier", func(t *testing.T) {
			identifier := exchangeClient.GetExchangeIdentifier()
			assert.NotEmpty(t, identifier)
		}) && pass
		assert.Equal(t, true, pass)
	}
}
