package coinbasepro

import (
	"strings"
)

func GetCoinbaseProSymbolFromCoinprice(
	coinpriceSymbol string,
) string {
	if symbol, ok := coinpriceToCoinbase[coinpriceSymbol]; ok {
		return symbol
	}
	return coinpriceSymbol
}
func GetCoinpriceSymbolFromCoinbasePro(
	coinbaseSymbol string,
) string {
	if symbol, ok := coinbaseToCoinprice[coinbaseSymbol]; ok {
		return symbol
	}
	return strings.ToUpper(coinbaseSymbol)
}
