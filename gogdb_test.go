package godbg

import (
	"fmt"
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
		Convey("Test custom buffer on global pdbg", func() {
			SetBuffers(nil)
			fmt.Fprintln(Out(), "test content")
			fmt.Fprintln(Err(), "err1 cerr")
			fmt.Fprintln(Err(), "err2 cerr2")
			fmt.Fprint(Out(), "test2 content2")
			So(OutString(), ShouldEqual, `test content
test2 content2`)
			So(ErrString(), ShouldEqual, `err1 cerr
err2 cerr2
`)
		})
	})
}
