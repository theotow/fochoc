package main

import (
	bit "github.com/toorop/go-bittrex"
)

type bittrex struct {
	ResultRaw []bit.Balance
}
type methodsBittrex struct{}

func initBittrex(c ConfigInterface) *bittrex {
	api := bit.New(c.GetKey("BITTREX_KEY"), c.GetKey("BITTREX_SECRET"))
	result, err := api.GetBalances()

	if err != nil {
		panic(err)
	}
	// val, _ := decimal.NewFromString("1.00")
	// result = append(result, bit.Balance{
	// 	Currency:      "ETH",
	// 	Balance:       val,
	// 	Available:     val,
	// 	Pending:       val,
	// 	CryptoAddress: "test",
	// 	Requested:     true,
	// 	Uuid:          "string",
	// })
	return &bittrex{ResultRaw: result}
}

func (m *methodsBittrex) Get(c ConfigInterface) ConfigProviderInterface {
	return initBittrex(c)
}

func (m *methodsBittrex) ConfigKeys() []string {
	return []string{"BITTREX_KEY", "BITTREX_SECRET"}
}

func NewBittrex() *methodsBittrex {
	return &methodsBittrex{}
}

func (b *bittrex) GetCurrencyValue(name string) float64 {
	for _, balance := range b.ResultRaw {
		if balance.Currency == name {
			val, _ := balance.Balance.Float64()
			return val
		}
	}
	return 0.00
}

func (b *bittrex) GetAll(keys []string) map[string]float64 {
	m := make(map[string]float64)
	for _, key := range keys {
		m[key] = b.GetCurrencyValue(key)
	}
	return m
}
