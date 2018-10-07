package main

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	EthAddress = "0xAd3cbAA752C8a5267785992Cf47723b992B33dD2"
	BtcAddress = "3MRT37Cja4TQtkL3cckC7AoNkX3dBEK39B"
)

var handler = func(input interface{}) (interface{}, error) {
	if b, ok := input.(BalanceSimple); ok {
		b.Comment = "comment"
		return b, nil
	}
	return nil, errors.New("handler error")
}

func TestScraper(t *testing.T) {
	initLogger(true)
	Convey("getER20Tokens() should execute without error", t, func() {
		res, _ := getER20Tokens("ETH", EthAddress)
		So(res["ETH"], ShouldEqual, 10430.564501815777)
		So(res["LRC"], ShouldEqual, 5)
		So(res["EMO"], ShouldEqual, 9063.343162)
		So(res["TRX"], ShouldEqual, 42)
		So(len(res), ShouldBeGreaterThan, 8)
	})
	Convey("getER20Tokens() should return empty map if address invalid", t, func() {
		res, _ := getER20Tokens("ETH", EthAddress+"32233")
		So(len(res), ShouldEqual, 0)
	})
	Convey("getCoins() should execute without error", t, func() {
		res := getCoins(1, 100)
		So(res["BTC"].Symbol, ShouldEqual, "BTC")
		So(res["BTC"].BtcPrice, ShouldEqual, 1)
		So(res["BTC"].UsdPrice, ShouldBeGreaterThan, 0)
	})
	Convey("getCoinsAsync() should execute without error", t, func() {
		res := getCoinsAsync()
		So(res["BTC"].Symbol, ShouldEqual, "BTC")
		So(res["BTC"].BtcPrice, ShouldEqual, 1)
		So(res["BTC"].UsdPrice, ShouldBeGreaterThan, 1000)
		So(len(res), ShouldBeGreaterThan, 1830)
	})
	Convey("getBalanceChainz() should execute without error", t, func() {
		res, _ := getBalanceChainz("BTC", BtcAddress)
		So(res["BTC"], ShouldEqual, 0)
	})
	Convey("getBalanceChainz() should execute have error", t, func() {
		_, err := getBalanceChainz("BTCCC", BtcAddress)
		So(err.Error(), ShouldEqual, "parse float blanace of BTCCC failed")
	})
	Convey("asyncRun() should do success / error case", t, func() {
		// success case
		res := asyncRun([]AsyncTask{
			{
				task: BalanceSimple{
					Balance: 11,
					Address: "address",
				},
				exec: handler,
			},
		}, 10)
		resFirst, ok := res[0].GetResult().(BalanceSimple)
		So(ok, ShouldEqual, true)
		So(len(res), ShouldEqual, 1)
		So(resFirst.Comment, ShouldEqual, "comment")
		So(resFirst.Balance, ShouldEqual, 11)
		So(resFirst.Address, ShouldEqual, "address")

		// error case
		res1 := asyncRun([]AsyncTask{
			{
				task: ProviderError{
					errorMsg:   "msg",
					providerID: "id",
				},
				exec: handler,
			},
		}, 10)
		So(res1[0].error.Error(), ShouldEqual, "handler error")

		// empty case
		res2 := asyncRun([]AsyncTask{}, 10)
		So(len(res2), ShouldEqual, 0)
	})
	Convey("resolveCoins() should run without error", t, func() {
		array := []BalanceSimple{
			{
				Address:  EthAddress,
				Currency: "ETH",
			},
			{
				Address:  BtcAddress,
				Currency: "BTC",
			},
		}
		mappedData := map[string]BalanceSimple{}
		res := resolveCoins(array)
		for _, val := range res {
			mappedData[val.Currency] = val
		}
		So(mappedData["ETH"].Balance, ShouldEqual, 10430.564501815777)
		So(mappedData["BTC"].Balance, ShouldEqual, 0)
	})
	Convey("getBalanceLisk() should get lisk balance", t, func() {
		res, err := getBalanceLisk("LSK", "6557304785210489363L")
		So(err, ShouldBeNil)
		So(res["LSK"], ShouldEqual, 0)
	})
	Convey("getBalanceEtc() should get etc balance", t, func() {
		res, err := getBalanceEtc("ETC", "0xb09ca4047ec095fb1dd6c8f916789056fed02615")
		So(err, ShouldBeNil)
		So(res["ETC"], ShouldEqual, 3016.3069)
	})
	Convey("getBalanceNeo() should get neo balance", t, func() {
		res, err := getBalanceNeo("NEO", "AaZiNiSmSHpZmUEK5PR7uXVSZWbFd5wfkb")
		So(err, ShouldBeNil)
		So(res["NEO"], ShouldEqual, 158)
	})
}
