package ftx

import "strings"

func GetCoinpriceSymbolFtx(
	coinbaseSymbol string,
) string {
	if symbol, ok := ftxToCoinprice[coinbaseSymbol]; ok {
		return symbol
	}
	return strings.ToUpper(coinbaseSymbol)
}
