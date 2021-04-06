package models

import (
	"time"
)

type LastSync struct {
	BaseAsset    string
	QuoteAsset   string
	Exchange     string
	LastSyncTime time.Time
}
