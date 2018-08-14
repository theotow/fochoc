package main

import (
	"errors"
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
	result["LTC"] = polo.Balance{
		Available: "5.015",
		BtcValue:  "0.078",
		OnOrders:  "0.078",
	}
	return &poloniex{ResultRaw: result}
}

func (m *methodsPoloniex) Get(c ConfigInterface) ConfigProviderInterface {
	return initPoloniex(c)
}

func (m *methodsPoloniex) ConfigKeys() []string {
	return []string{"POLONIEX_KEY", "POLONIEX_SECRET"}
}

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

func (b *poloniex) GetAll(keys []string) map[string]float64 {
	m := make(map[string]float64)
	for _, key := range keys {
		m[key] = b.GetCurrencyValue(key)
	}
	return m
}
