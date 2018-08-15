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
	SkipConvey("should be able to get key", t, func() {

	})
	SkipConvey("should get all", t, func() {

	})
}
