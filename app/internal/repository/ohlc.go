package repository

import (
	"github.com/chain-bot/prices/app/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/volatiletech/null/v8"
	"time"
)

func (repo *RepositoryImpl) UpsertOHLCVData(
	ohlcvMarketData []*models.OHLCVMarketData,
	exchange string,
	pair *models.Symbol,
) {
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
		repo.writeAPI.WritePoint(p)
	}
}

func (repo *RepositoryImpl) GetOHLCVData(
	base string,
	quote null.String,
	exchange null.String,
	start time.Time,
	end time.Time,
) ([]models.OHLCVMarketData, error) {
	return []models.OHLCVMarketData{}, nil
}
