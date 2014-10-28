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
			pdbg.bout = nil
			pdbg.sout = nil
			pdbg.berr = nil
			pdbg.serr = nil
			fmt.Fprintln(Out(), "test0 content0")
			So(OutString(), ShouldEqual, ``)
			fmt.Fprintln(Err(), "err0 content0")
			So(ErrString(), ShouldEqual, ``)
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

		Convey("Test custom buffer reset on global pdbg", func() {
			SetBuffers(nil)
			fmt.Fprint(Out(), "test content")
			So(OutString(), ShouldEqual, `test content`)
			fmt.Fprint(Err(), "err1 cerr")
			So(ErrString(), ShouldEqual, `err1 cerr`)
			ResetIOs()
			fmt.Fprint(Out(), "test2 content2")
			So(OutString(), ShouldEqual, `test2 content2`)
			fmt.Fprint(Err(), "err2 cerr2")
			So(ErrString(), ShouldEqual, `err2 cerr2`)
		})

		Convey("Test custom buffer on custom pdbg", func() {
			apdbg := NewPdbg(SetBuffers)
			fmt.Fprintln(apdbg.Out(), "test content")
			fmt.Fprintln(apdbg.Err(), "err1 cerr")
			fmt.Fprintln(apdbg.Err(), "err2 cerr2")
			fmt.Fprint(apdbg.Out(), "test2 content2")
			So(apdbg.OutString(), ShouldEqual, `test content
test2 content2`)
			So(apdbg.ErrString(), ShouldEqual, `err1 cerr
err2 cerr2
`)
		})
		Convey("Test custom buffer reset on custom pdbg", func() {
			apdbg := NewPdbg(SetBuffers)
			fmt.Fprint(apdbg.Out(), "test content")
			So(apdbg.OutString(), ShouldEqual, `test content`)
			fmt.Fprint(apdbg.Err(), "err1 cerr")
			So(apdbg.ErrString(), ShouldEqual, `err1 cerr`)
			apdbg.ResetIOs()
			fmt.Fprint(apdbg.Out(), "test2 content2")
			So(apdbg.OutString(), ShouldEqual, `test2 content2`)
			fmt.Fprint(apdbg.Err(), "err2 cerr2")
			So(apdbg.ErrString(), ShouldEqual, `err2 cerr2`)
		})
	})

	Convey("Test pdbg print functions", t, func() {
		Convey("Test pdbg print with global instance", func() {
			SetBuffers(nil)
			Pdbgf("test")
			So(ErrString(), ShouldEqual,
				`[func.010:95]
  test
`)
			ResetIOs()
			prbgtest()
			So(ErrString(), ShouldEqual,
				` [prbgtest:126] (func.010:101)
   prbgtest content
`)
		})

		Convey("Test pdbg print with custom instance", func() {
			apdbg := NewPdbg(SetBuffers)
			apdbg.Pdbgf("test2")
			So(apdbg.ErrString(), ShouldEqual,
				`[func.011:110]
  test2
`)
			apdbg.ResetIOs()
			prbgtestCustom(apdbg)
			So(apdbg.ErrString(), ShouldEqual,
				` [prbgtestCustom:130] (func.011:116)
   prbgtest content2
`)
		})
	})
}

func prbgtest() {
	Pdbgf("prbgtest content")
}

func prbgtestCustom(pdbg *Pdbg) {
	pdbg.Pdbgf("prbgtest content2")
}
