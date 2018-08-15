package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var pConfig = EnvConfig{}
var pInstance = NewPoloniex()

func TestPoloniex(t *testing.T) {
	Convey("should not crash", t, func() {
		pInstance.Get(&pConfig)
	})
	Convey("should be able to get key", t, func() {
		i := pInstance.Get(&pConfig)
		i.AddTestBalance("ETH", 0.1)
		result := i.GetCurrencyValue("ETH")
		So(result, ShouldEqual, 0.1)
		// does not exist
		notexisting := i.GetCurrencyValue("ETHH")
		So(notexisting, ShouldEqual, 0.0)
	})
	SkipConvey("should get all", t, func() {

	})
}
