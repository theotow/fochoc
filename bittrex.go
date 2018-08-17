package main

import (
	"github.com/shopspring/decimal"
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

func (b *bittrex) AddTestBalance(name string, value float64) {
	val := decimal.NewFromFloat(value)
	b.ResultRaw = append(b.ResultRaw, bit.Balance{
		Currency:      name,
		Balance:       val,
		Available:     val,
		Pending:       val,
		CryptoAddress: "test",
		Requested:     true,
		Uuid:          "string",
	})
}

func (b *bittrex) GetAll(keys []string) []BalanceSimple {
	return GetAllValues(keys, b.GetCurrencyValue)
}
