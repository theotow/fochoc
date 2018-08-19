package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var pConfig = EnvConfig{}
var pInstance = NewPoloniex()

func TestPoloniex(t *testing.T) {
	Convey("should not crash", t, func() {
		_, err := pInstance.Get(&pConfig)
		So(err, ShouldBeNil)
	})
	Convey("should be able to get key", t, func() {
		i, _ := pInstance.Get(&pConfig)
		i.AddTestBalance("ETH", 0.1)
		result := i.GetCurrencyValue("ETH")
		So(result, ShouldEqual, 0.1)
		// does not exist
		notexisting := i.GetCurrencyValue("ETHH")
		So(notexisting, ShouldEqual, 0.0)
	})
	Convey("should be able to get all", t, func() {
		i, _ := pInstance.Get(&pConfig)
		i.AddTestBalance("ETH", 1.11)
		res := i.GetAll([]string{"ETH"})
		So(len(res), ShouldEqual, 1)
		So(res[0].Balance, ShouldEqual, 1.11)
	})
}
