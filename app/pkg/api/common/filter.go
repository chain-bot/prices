package common

import (
	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/pkg/models"
)

func FilterSupportedAssets(symbols []*models.Symbol) []*models.Symbol {
	var result []*models.Symbol
	supportedAssets := configs.GetSupportedAssets()
	for index := range symbols {
		pair := symbols[index]
		_, ok := supportedAssets[pair.NormalizedBase]
		if !ok {
			continue
		}
		_, ok = supportedAssets[pair.NormalizedQuote]
		if !ok {
			continue
		}
		result = append(result, pair)
	}
	return result
}
