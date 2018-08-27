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
	}
	return ""
}

func (c Config) GetColdWalletCoins() []ColdWalletCoin {
	return []ColdWalletCoin{}
}

func TestRun(t *testing.T) {
	Convey("initProviders() should only init valid providers", t, func() {
		config := Config{}
		res, _ := initProviders([]string{"binance", "kraken"}, config)
		So(len(res), ShouldEqual, 1)
		So(res[0].id, ShouldEqual, "kraken")
	})
	Convey("getAggSum() should return btc / usd sum", t, func() {
		input := []Balance{
			NewBalance("ETH", 1, 1, 0.5, 3000, "binance", NewBinance()),
			NewBalance("BTC", 2, 1, 1, 6000, "binance", NewBinance()),
		}
		usd, btc := getAggSum(input)
		So(usd, ShouldEqual, 15000)
		So(btc, ShouldEqual, 2.5)
	})
	Convey("getCoinsOfProvider() should return coins of provider", t, func() {
		config := EnvConfig{}
		instance, err := NewBinance().Get(&config)
		So(err, ShouldBeNil)
		instance.AddTestBalance("BTC", 2) // should not show up
		instance.AddTestBalance("ETH", 2) // should show up
		provider := Provider{
			instance: instance,
			factory:  NewBinance(),
			id:       "binance",
		}
		mapCoins := make(map[string]Coin)
		mapCoins["ETH"] = Coin{
			Symbol: "ETH",
		} // add coin to known coins list
		balances := provider.getCoinsOfProvider(mapCoins)
		So(len(balances), ShouldEqual, 1)
		So(balances[0].Balance.Balance, ShouldEqual, 2)
		So(balances[0].Balance.Currency, ShouldEqual, "ETH")
		So(balances[0].Coin.Symbol, ShouldEqual, "ETH")
		So(balances[0].Provider.id, ShouldEqual, "binance")
	})
}
