package api

import "time"

type Granularity string

const (
	MinuteGranularity Granularity = "min"
	HourGranularity   Granularity = "hour"
	DayGranularity    Granularity = "day"
)

type ExchangeAPIClient interface {
	GetPrices(startTime time.Time, endTime time.Time, granularity Granularity) ([]*OHLCPriceVolume, error)
}

type OHLCPriceVolume struct {
	open   float64
	high   float64
	low    float64
	close  float64
	volume float64
}
