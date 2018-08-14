package config

import (
	"errors"
	"os"
	"testing"

	"io/ioutil"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRun(t *testing.T) {
	Convey("should get file url", t, func() {
		res := getFileString()
		if len(res) == 0 {
			t.Error("should have len > 0")
		}
	})
	Convey("should panic if error", t, func() {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("should have panicked!")
				}
			}()
			check(errors.New("error"), "errormsg")
		}()
	})
	Convey("should read / write config", t, func() {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			panic(err)
		}
		path := tmpfile.Name()
		defer os.Remove(path)
		writeFile(path, map[string]string{
			"a": "b",
		})
		data := readFile(path)
		writeFile(path, data)
		So(data, ShouldContainKey, "a")
		So(data["a"], ShouldEqual, "b")

	})
}
