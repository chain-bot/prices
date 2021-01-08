package binance

import "time"

func UnixMillis(t time.Time) int64 {
	t.UTC()
	return t.UnixNano() / int64(time.Millisecond)
}
