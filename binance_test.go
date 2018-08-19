package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var bConfig = EnvConfig{}
var bInstance = NewBinance()

func TestBinance(t *testing.T) {
	Convey("should not crash", t, func() {
		_, err := bInstance.Get(&bConfig)
		So(err, ShouldBeNil)
	})
	Convey("should be able to get key", t, func() {
		res, err := bInstance.Get(&bConfig)
		So(err, ShouldBeNil)
		res.AddTestBalance("ETH", 1.01)
		result := res.GetCurrencyValue("ETH")
		So(result, ShouldEqual, 1.01)
	})
	Convey("should be able to get all", t, func() {
		i, err := bInstance.Get(&bConfig)
		So(err, ShouldBeNil)
		i.AddTestBalance("ETH", 1.11)
		res := i.GetAll([]string{"ETH"})
		So(len(res), ShouldEqual, 1)
		So(res[0].Balance, ShouldEqual, 1.11)
	})
}
