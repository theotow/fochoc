package binance_test

import (
	"go-crypto/binance"
	configImport "go-crypto/config"
	"testing"

	bin "github.com/adshao/go-binance"
)

var config = configImport.EnvConfig{}

func TestInit(t *testing.T) {
	res := binance.Get(&config)
	if len(res.ResultRaw) != 0 {
		t.Error("should be empty array")
	}
}

func TestGetKey(t *testing.T) {
	res := binance.Get(&config)
	res.ResultRaw = []bin.Balance{
		bin.Balance{Asset: "ETH", Free: "0.01", Locked: "1.00"},
	}
	result := res.GetCurrencyValue("ETH")
	if result != 1.01 {
		t.Error("should have 1.01 balance for ETH")
	}
}

func TestGetKeyError(t *testing.T) {
	res := binance.Get(&config)
	res.ResultRaw = []bin.Balance{
		bin.Balance{Asset: "ETH", Free: "invalid", Locked: "invalid"},
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("should panic if invalid float")
			}
		}()
		res.GetCurrencyValue("ETH")
	}()
}

func TestAll(t *testing.T) {
	res := binance.Get(&config)
	res.ResultRaw = []bin.Balance{
		bin.Balance{Asset: "ETH", Free: "0.00", Locked: "1.00"},
		bin.Balance{Asset: "XMR", Free: "0.00", Locked: "1.00"},
	}
	result := res.GetAll([]string{"ETH", "XMR"})
	if len(result) != 2 {
		t.Error("should have len == 2")
	}
	for _, key := range []string{"ETH", "XMR"} {
		if result[key] != 1 {
			t.Error("should be 1")
		}
	}
}
