package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var bConfig = EnvConfig{}
var bInstance = NewBinance()

func TestBinance(t *testing.T) {
	Convey("should not crash", t, func() {
		bInstance.Get(&bConfig)
	})
	Convey("should be able to get key", t, func() {
		res := bInstance.Get(&bConfig)
		res.AddTestBalance("ETH", 1.01)
		result := res.GetCurrencyValue("ETH")
		So(result, ShouldEqual, 1.01)
	})
	SkipConvey("should get all", t, func() {
		// res := bInstance.Get(&bConfig)
		// res.ResultRaw = []bin.Balance{
		// 	bin.Balance{Asset: "ETH", Free: "0.00", Locked: "1.00"},
		// 	bin.Balance{Asset: "XMR", Free: "0.00", Locked: "1.00"},
		// }
		// result := res.GetAll([]string{"ETH", "XMR"})
		// if len(result) != 2 {
		// 	t.Error("should have len == 2")
		// }
		// for _, key := range []string{"ETH", "XMR"} {
		// 	if result[key] != 1 {
		// 		t.Error("should be 1")
		// 	}
		// }
	})
}
