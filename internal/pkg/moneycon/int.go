package moneycon

func FloatToInt(f float64) int64 {
	return int64(f * 100)
}

func IntToFloat(i int64) float64 {
	return float64(i) / 100
}
