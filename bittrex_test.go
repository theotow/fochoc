package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var biConfig = EnvConfig{}
var biInstance = NewBittrex()

func TestBittrex(t *testing.T) {
	Convey("should not crash", t, func() {
		biInstance.Get(&biConfig)
	})
	Convey("should be able to get key", t, func() {
		i := biInstance.Get(&biConfig)
		i.AddTestBalance("ETH", 1.11)
		So(i.GetCurrencyValue("ETH"), ShouldEqual, 1.11)
	})
	Convey("should be able to get all", t, func() {
		i := biInstance.Get(&biConfig)
		i.AddTestBalance("ETH", 1.11)
		res := i.GetAll([]string{"ETH"})
		So(len(res), ShouldEqual, 1)
		So(res[0].Balance, ShouldEqual, 1.11)
	})
}
