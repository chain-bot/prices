package common

import (
	"github.com/mochahub/coinprice-scraper/app/configs"
	"github.com/mochahub/coinprice-scraper/app/pkg/models"
)

func FilterSupportedAssets(symbols []*models.Symbol) []*models.Symbol {
	result := []*models.Symbol{}
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
