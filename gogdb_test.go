package godbg

import (
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProject(t *testing.T) {
	Convey("Test buffers", t, func() {

		Convey("By Default, equals to std", func() {
			So(Out(), ShouldEqual, os.Stdout)
			So(Err(), ShouldEqual, os.Stderr)
		})
		Convey("When set to buffer, no longer equals to std", func() {
			SetBuffers(nil)
			So(Out(), ShouldNotEqual, os.Stdout)
			So(Err(), ShouldNotEqual, os.Stderr)
		})
	})
}
