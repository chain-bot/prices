package common

import (
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/scraper/models"
)

func FilterSupportedAssets(symbols []*models.Symbol) []*models.Symbol {
	result := []*models.Symbol{}
	supportedAssets := config.GetSupportedAssets()
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
