package main

import (
	"errors"
	"fmt"
	"strconv"

	polo "github.com/jyap808/go-poloniex"
)

type poloniex struct {
	ResultRaw map[string]polo.Balance
}
type methodsPoloniex struct{}

func initPoloniex(c ConfigInterface) *poloniex {
	api := polo.New(c.GetKey("POLONIEX_KEY"), c.GetKey("POLONIEX_SECRET"))
	result, err := api.GetBalances()

	if err != nil {
		panic(err)
	}
	return &poloniex{ResultRaw: result}
}

func (m *methodsPoloniex) Get(c ConfigInterface) ConfigProviderInterface {
	return initPoloniex(c)
}

func (m *methodsPoloniex) ConfigKeys() []string {
	return []string{"POLONIEX_KEY", "POLONIEX_SECRET"}
}

// NewPoloniex is used to create a poloniex provider adapter
func NewPoloniex() *methodsPoloniex {
	return &methodsPoloniex{}
}

func (b *poloniex) GetCurrencyValue(name string) float64 {
	for currency, balance := range b.ResultRaw {
		if currency == name {
			free, errFree := strconv.ParseFloat(balance.Available, 64)
			locked, errLocked := strconv.ParseFloat(balance.OnOrders, 64)
			if errFree != nil || errLocked != nil {
				panic(errors.New("could not parse int"))
			}
			return free + locked
		}
	}
	return 0.00
}

func (b *poloniex) AddTestBalance(name string, value float64) {
	b.ResultRaw[name] = polo.Balance{
		Available: fmt.Sprintf("%f", value),
		BtcValue:  "0.000",
		OnOrders:  "0.000",
	}
}

func (b *poloniex) GetAll(keys []string) []BalanceSimple {
	return GetAllValues(keys, b.GetCurrencyValue)
}
