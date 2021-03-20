package binance

import "time"

func getBinanceIntervalFromDuration(
	interval time.Duration,
) Interval {
	ret := Minute
	if int(interval.Hours()) > 0 {
		ret = Hour
	}
	return ret
}
