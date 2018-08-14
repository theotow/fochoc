package kraken

import (
	"go-crypto/config"
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

type methods struct{}

func Get(c config.ConfigInterface) *kraken {
	api := krakenapi.New(c.GetKey("KRAKEN_KEY"), c.GetKey("KRAKEN_SECRET"))
	result, err := api.Balance()

	if err != nil {
		panic(err)
	}
	// result.XETH = 10.00
	return &kraken{ResultRaw: result}
}

func (m *methods) Get(c config.ConfigInterface) config.ProviderInterface {
	return Get(c)
}

func (m *methods) ConfigKeys() []string {
	return []string{"KRAKEN_KEY", "KRAKEN_SECRET"}
}

func New() *methods {
	return &methods{}
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
