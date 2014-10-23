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
	})
}
