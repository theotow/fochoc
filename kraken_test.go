package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var config = EnvConfig{}
var instance = NewKraken()

func TestKraken(t *testing.T) {
	Convey("should not crash", t, func() {
		instance.Get(&config)
	})
	SkipConvey("should be able to get key", t, func() {
		res := instance.Get(&config)
		// res.ResultRaw.XETH = 0.1
		result := res.GetCurrencyValue("ETH")
		if result != 0.1 {
			t.Error("Should have 0.1 balance for ETH")
		}
		// does not exist
		notexisting := res.GetCurrencyValue("ETHH")
		if notexisting != 0 {
			t.Error("Should default to 0")
		}
	})
	SkipConvey("should get all", t, func() {
		res := instance.Get(&config)
		// res.ResultRaw.XETH = 1
		// res.ResultRaw.XXMR = 1
		result := res.GetAll([]string{"ETH", "XMR"})
		if len(result) != 2 {
			t.Error("should have len == 2")
		}
		for _, key := range []string{"ETH", "XMR"} {
			if result[key] != 1 {
				t.Error("should be 1")
			}
		}
	})
}
