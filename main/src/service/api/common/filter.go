package common

import "github.com/mochahub/coinprice-scraper/config"

func FilterSupportedAssets(symbols []*Symbol) []*Symbol {
	result := []*Symbol{}
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
