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
	SkipConvey("should be able to get key", t, func() {

	})
	SkipConvey("should get all", t, func() {

	})
}
