package binance

import (
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api/common"
	"github.com/mochahub/coinprice-scraper/app/utils"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestBinanceClient(t *testing.T) {
	// TODO: Use DI instead of calling GetSecrets directly
	utils.LoadEnv()
	secret, _ := config.GetSecrets()
	exchangeClient := NewBinanceAPIClient(secret)
	pass := true
	pass = t.Run("TestGetCandleStickData", func(t *testing.T) {
		expectedLength := 480 * time.Minute
		startTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.Add(expectedLength - time.Minute)
		candleStickData, err := exchangeClient.GetOHLCMarketData(
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
	pass = t.Run("TestGetExchangeInfo", func(t *testing.T) {
		exchangeInfo, err := exchangeClient.getExchangeInfo()
		assert.NoError(t, err)
		assert.NotNil(t, exchangeInfo)
		//fmt.Print(utils.PrettyJSON(exchangeInfo))
	}) && pass
	// Interface Methods
	// TODO(Zahin): Do we even need this? exhange_clients_test will test it as well...
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
		log.Println(candleStickData[len(candleStickData)-1].StartTime.String())
		log.Println(candleStickData[0].StartTime.String())
		assert.NoError(t, err)
		assert.NotEmpty(t, candleStickData)
		assert.Equal(t, int(expectedLength.Minutes()), len(candleStickData))
	}) && pass

	pass = t.Run("TestGetSupportedPairs", func(t *testing.T) {
		pairs, err := exchangeClient.GetSupportedPairs()
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
		//fmt.Print(utils.PrettyJSON(pairs))
		assert.Equal(t, 3, len(pairs))
	}) && pass
	assert.Equal(t, true, pass)
}
