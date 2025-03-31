package token

import "math"

func ConvertToLamport(amount float64) uint64 {
	return uint64(math.Round(amount * math.Pow10(9)))
}

func ConvertToDecimals(human float64, decimals uint8) uint64 {
	factor := math.Pow10(int(decimals))
	rp := human * factor
	return uint64(math.Round(rp))
}
