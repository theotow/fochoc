package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var config = EnvConfig{}
var instance = NewKraken()

func TestKraken(t *testing.T) {
	Convey("should not crash", t, func() {
		_, err := instance.Get(&config)
		So(err, ShouldBeNil)
	})
	Convey("should be able to get key", t, func() {
		res, _ := instance.Get(&config)
		res.AddTestBalance("ETH", 0.1)
		result := res.GetCurrencyValue("ETH")
		So(result, ShouldEqual, 0.1)
		// does not exist
		notexisting := res.GetCurrencyValue("ETHH")
		So(notexisting, ShouldEqual, 0.0)
	})
	Convey("should be able to get all", t, func() {
		i, _ := instance.Get(&config)
		i.AddTestBalance("ETH", 1.11)
		res := i.GetAll([]string{"ETH"})
		So(len(res), ShouldEqual, 1)
		So(res[0].Balance, ShouldEqual, 1.11)
	})
}
