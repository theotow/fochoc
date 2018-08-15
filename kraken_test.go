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
	Convey("should be able to get key", t, func() {
		res := instance.Get(&config)
		res.AddTestBalance("ETH", 0.1)
		result := res.GetCurrencyValue("ETH")
		So(result, ShouldEqual, 0.1)
		// does not exist
		notexisting := res.GetCurrencyValue("ETHH")
		So(notexisting, ShouldEqual, 0.0)
	})
	SkipConvey("should get all", t, func() {

	})
}
