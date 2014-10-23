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
		Convey("By Default, a new pdbg instance buffer equals to std", func() {
			apdbg := NewPdbg()
			So(apdbg.Out(), ShouldEqual, os.Stdout)
			So(apdbg.Err(), ShouldEqual, os.Stderr)
		})
		Convey("By Default, a new pdbg instance set to buffer writes no longer equals to std", func() {
			apdbg := NewPdbg(SetBuffers)
			So(apdbg.Out(), ShouldNotEqual, os.Stdout)
			So(apdbg.Err(), ShouldNotEqual, os.Stderr)
		})
	})
}
