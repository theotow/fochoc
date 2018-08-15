package main

import (
	"errors"
	"os"
	"testing"

	"io/ioutil"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
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
		input := ConfigFileStruct{
			Keys: map[string]string{
				"a": "b",
			},
		}
		writeFile(path, input)
		output := readFile(path)
		writeFile(path, output)
		So(output.Keys, ShouldContainKey, "a")
		So(output.Keys["a"], ShouldEqual, "b")
	})
}
