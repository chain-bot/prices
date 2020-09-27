package api

type Interval string

const (
	BinanceMinuteInterval Interval = "1m"
	BinanceHourInterval   Interval = "1h"
	BinanceDayInterval    Interval = "1D"
)
