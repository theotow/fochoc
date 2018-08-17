package main

func GetAllValues(keys []string, getter func(string) float64) []BalanceSimple {
	output := []BalanceSimple{}
	for _, key := range keys {
		output = append(output, BalanceSimple{
			Currency: key,
			Balance:  getter(key),
			Comment:  "-",
			Address:  "-",
		})
	}
	return output
}
