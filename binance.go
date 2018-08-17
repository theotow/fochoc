package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	bin "github.com/adshao/go-binance"
)

type binance struct {
	ResultRaw []bin.Balance
}
type methodsBinance struct{}

func initBinance(c ConfigInterface) *binance {
	api := bin.NewClient(c.GetKey("BINANCE_KEY"), c.GetKey("BINANCE_SECRET"))
	result, err := api.NewGetAccountService().Do(context.Background())

	if err != nil {
		panic(err)
	}
	return &binance{ResultRaw: result.Balances}
}

func (m *methodsBinance) Get(c ConfigInterface) ConfigProviderInterface {
	return initBinance(c)
}

func (m *methodsBinance) ConfigKeys() []string {
	return []string{"BINANCE_KEY", "BINANCE_SECRET"}
}

// NewBinance is used to create a binance provider adapter
func NewBinance() *methodsBinance {
	return &methodsBinance{}
}

func (b *binance) GetCurrencyValue(name string) float64 {
	for _, balance := range b.ResultRaw {
		if balance.Asset == name {
			free, errFree := strconv.ParseFloat(balance.Free, 64)
			locked, errLocked := strconv.ParseFloat(balance.Locked, 64)
			if errFree != nil || errLocked != nil {
				panic(errors.New("could not parse int"))
			}
			return free + locked
		}
	}
	return 0.00
}

func (b *binance) AddTestBalance(name string, value float64) {
	b.ResultRaw = append(b.ResultRaw, bin.Balance{Free: fmt.Sprintf("%f", value), Locked: "0.00", Asset: name})
}

func (b *binance) GetAll(keys []string) []BalanceSimple {
	return GetAllValues(keys, b.GetCurrencyValue)
}
