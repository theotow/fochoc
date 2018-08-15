package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var biConfig = EnvConfig{}
var biInstance = NewBinance()

func TestBittrex(t *testing.T) {
	Convey("should not crash", t, func() {
		biInstance.Get(&biConfig)
	})
	Convey("should be able to get key", t, func() {
		i := biInstance.Get(&biConfig)
		i.AddTestBalance("ETH", 1.11)
		So(i.GetCurrencyValue("ETH"), ShouldEqual, 1.11)
	})
	SkipConvey("should get all", t, func() {

	})
}
