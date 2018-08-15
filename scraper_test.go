package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScraper(t *testing.T) {
	Convey("getER20Tokens() should execute without error", t, func() {
		res := GetER20Tokens("0xf27d22d64e625c2a34e31369d9b88828146df52b")
		So(res["ETH"], ShouldEqual, 0.101208623)
		So(res["DAR"], ShouldEqual, 2)
		So(len(res), ShouldEqual, 3)
	})
}
