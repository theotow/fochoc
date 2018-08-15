package main

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Config struct{}

func (c Config) GetKey(name string) string {
	if name == "KRAKEN_KEY" || name == "KRAKEN_SECRET" {
		return os.Getenv(name)
	} else {
		return ""
	}
}

func (c Config) Initialised() bool {
	return true
}

func (c Config) GetTokens() []Token {
	return []Token{}
}

func TestRun(t *testing.T) {
	Convey("getCoins() should execute without error", t, func() {
		res := getCoins()
		if res["BTC"].Symbol != "BTC" {
			t.Error("btc not found")
		}
		if res["BTC"].BtcPrice != 1 {
			t.Error("btc price should be 1")
		}
		if res["BTC"].UsdPrice <= 0 {
			t.Error("usd price should be > 0")
		}
	})
	Convey("initProviders() should only init valid providers", t, func() {
		config := Config{}
		res := initProviders([]string{"binance", "kraken"}, config)
		So(len(res), ShouldEqual, 1)
		So(res[0].id, ShouldEqual, "kraken")
	})
}
