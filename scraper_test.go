package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScraper(t *testing.T) {
	Convey("getER20Tokens() should execute without error", t, func() {
		res := GetER20Tokens("0xAd3cbAA752C8a5267785992Cf47723b992B33dD2")
		So(res["ETH"], ShouldEqual, 10430.564501815777)
		So(res["LRC"], ShouldEqual, 5)
		So(res["EMO"], ShouldEqual, 9063.343162)
		So(res["TRX"], ShouldEqual, 42)
		So(len(res), ShouldBeGreaterThan, 8)
	})
}
