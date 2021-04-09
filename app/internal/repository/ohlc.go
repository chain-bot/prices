package repository

import (
	"github.com/chain-bot/scraper/app/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func (repo *RepositoryImpl) UpsertOHLCVData(
	ohlcvMarketData []*models.OHLCVMarketData,
	exchange string,
	pair *models.Symbol,
) {
	writeAPI := (*repo.influxClient).WriteAPI(repo.influxOrg, repo.ohlcvBucket)
	tags := map[string]string{
		"quote":    pair.NormalizedQuote,
		"exchange": exchange,
	}
	for index := range ohlcvMarketData {
		ohlcv := ohlcvMarketData[index]
		fields := map[string]interface{}{
			"open":   ohlcv.OpenPrice,
			"high":   ohlcv.HighPrice,
			"low":    ohlcv.LowPrice,
			"close":  ohlcv.ClosePrice,
			"volume": ohlcv.Volume,
		}
		p := influxdb2.NewPoint(
			pair.NormalizedBase,
			tags,
			fields,
			ohlcv.StartTime)
		writeAPI.WritePoint(p)
	}
}
