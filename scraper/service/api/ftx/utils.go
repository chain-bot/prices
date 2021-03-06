package ftx

import "strings"

func GetFtxSymbolFromCoinprice(
	coinpriceSymbol string,
) string {
	if symbol, ok := coinpriceToftx[coinpriceSymbol]; ok {
		return symbol
	}
	return coinpriceSymbol
}
func GetCoinpriceSymbolFtx(
	coinbaseSymbol string,
) string {
	if symbol, ok := ftxToCoinprice[coinbaseSymbol]; ok {
		return symbol
	}
	return strings.ToUpper(coinbaseSymbol)
}
