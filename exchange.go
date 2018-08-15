package main

type Exchange struct{}

func (b *Exchange) GetCurrencyValue(name string) float64 {
	return 0.00
}

func (b *Exchange) GetAll(keys []string) []BalanceSimple {
	output := []BalanceSimple{}
	for _, key := range keys {
		output = append(output, BalanceSimple{
			Currency: key,
			Balance:  b.GetCurrencyValue(key),
			Comment:  "-",
			Address:  "-",
		})
	}
	return output
}
