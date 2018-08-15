package main

import (
	"reflect"

	"github.com/beldur/kraken-go-api-client"
)

type kraken struct {
	ResultRaw *krakenapi.BalanceResponse
}

var mapping = map[string]string{
	"ETH":  "XETH",
	"BCH":  "BCH",
	"DASH": "DASH",
	"USDT": "USDT",
	"XMR":  "XXMR",
	"EUR":  "ZEUR",
}

type methodsKraken struct{}

func initKraken(c ConfigInterface) *kraken {
	api := krakenapi.New(c.GetKey("KRAKEN_KEY"), c.GetKey("KRAKEN_SECRET"))
	result, err := api.Balance()
	if err != nil {
		panic(err)
	}
	return &kraken{ResultRaw: result}
}

func (m *methodsKraken) Get(c ConfigInterface) ConfigProviderInterface {
	return initKraken(c)
}

func (m *methodsKraken) ConfigKeys() []string {
	return []string{"KRAKEN_KEY", "KRAKEN_SECRET"}
}

func NewKraken() *methodsKraken {
	return &methodsKraken{}
}

func (k *kraken) GetCurrencyValue(name string) float64 {
	r := reflect.ValueOf(k.ResultRaw)
	f := reflect.Indirect(r).FieldByName(k.getLocalKey(name))
	if f.IsValid() != true {
		return 0.00
	}
	return f.Float()
}

func (k *kraken) getLocalKey(key string) string {
	if val, ok := mapping[key]; ok {
		return val
	}
	// log (missing static mapping )
	return key
}

func (k *kraken) GetAll(keys []string) map[string]float64 {
	m := make(map[string]float64)
	for _, key := range keys {
		m[key] = k.GetCurrencyValue(key)
	}
	return m
}

func (k *kraken) AddTestBalance(name string, value float64) {
	v := reflect.ValueOf(k.ResultRaw).Elem().FieldByName(k.getLocalKey(name))
	if v.IsValid() {
		v.SetFloat(value)
	}
}
