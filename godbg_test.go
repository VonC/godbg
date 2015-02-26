package godbg

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProject(t *testing.T) {
	Convey("Test buffers", t, func() {

		Convey("By Default, equals to std", func() {
			So(Out(), ShouldEqual, os.Stdout)
			So(Err(), ShouldEqual, os.Stderr)
			So(NoOutput(), ShouldBeTrue)
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
			So(apdbg.NoOutput(), ShouldBeTrue)
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
			So(NoOutput(), ShouldBeFalse)
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
			So(apdbg.NoOutput(), ShouldBeFalse)
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
			So(ErrString(), ShouldEqualNL,
				`[func.012:96]
  test`)
			ResetIOs()
			prbgtest()
			So(ErrString(), ShouldEqualNL,
				`  [prbgtest:4] (func.012:102)
    prbgtest content`)
		})

		Convey("Test pdbg print with custom instance", func() {
			apdbg := NewPdbg(SetBuffers)
			apdbg.Pdbgf("test2")
			So(apdbg.ErrString(), ShouldEqualNL,
				`[func.013:111]
  test2`)
			apdbg.ResetIOs()
			prbgtestCustom(apdbg)
			So(apdbg.ErrString(), ShouldEqualNL,
				`  [prbgtestCustom:8] (func.013:117)
    prbgtest content2`)
			apdbg.ResetIOs()
			apdbg.pdbgTestInstance()
			So(apdbg.ErrString(), ShouldEqualNL,
				`  [*Pdbg.pdbgTestInstance:12] (func.013:123)
    pdbgTestInstance content3`)
		})
		Convey("Test pdbg prints nothing if runtime.Caller fails", func() {
			mycaller = failCaller
			apdbg := NewPdbg(SetBuffers)
			apdbg.Pdbgf("test fail")
			So(apdbg.ErrString(), ShouldEqual, `  test fail
`)
			mycaller = runtime.Caller
		})
	})

	Convey("Test pdbg excludes functions", t, func() {
		Convey("Test pdbg exclude with global instance", func() {
			SetBuffers(nil)
			pdbg.SetExcludes([]string{"globalNo"})
			globalPdbgExcludeTest()
			So(ErrString(), ShouldEqualNL,
				`  [globalPdbgExcludeTest:16] (func.016:143)
    calling no
      [globalCNo:26] (globalPdbgExcludeTest:17) (func.016:143)
        gcalled2`)
		})
		Convey("Test pdbg exclude with custom instance", func() {
			apdbg := NewPdbg(SetBuffers, OptExcludes([]string{"customNo"}))
			customPdbgExcludeTest(apdbg)
			So(apdbg.ErrString(), ShouldEqualNL,
				`  [customPdbgExcludeTest:30] (func.017:153)
    calling cno
      [customCNo:40] (customPdbgExcludeTest:31) (func.017:153)
        ccalled2`)
		})
	})

	Convey("Test pdbg skips functions", t, func() {
		Convey("Test pdbg skip with global instance", func() {
			SetBuffers(nil)
			pdbg.SetSkips([]string{"globalNo"})
			globalPdbgExcludeTest()
			So(ErrString(), ShouldEqualNL,
				`  [globalPdbgExcludeTest:16] (func.019:167)
    calling no
      [globalCNo:26] (globalPdbgExcludeTest:17) (func.019:167)
        gcalled2`)
		})
		Convey("Test pdbg skip with custom instance", func() {
			apdbg := NewPdbg(SetBuffers, OptSkips([]string{"customNo"}))
			customPdbgExcludeTest(apdbg)
			So(ErrString(), ShouldEqualNL,
				`  [globalPdbgExcludeTest:16] (func.019:167)
    calling no
      [globalCNo:26] (globalPdbgExcludeTest:17) (func.019:167)
        gcalled2`)
		})
	})

	Convey("Test pdbg can ignore line number", t, func() {
		Convey("Test pdbg skip with global instance", func() {
			SetBuffers(nil)
			pdbg.SetSkips([]string{"globalNo"})
			globalPdbgExcludeTest()
			s := ShouldEqualNL(ErrString(), `r`)
			So(s, ShouldNotEqual, `e`)
		})
		Convey("Test pdbg skip with custom instance", func() {
			apdbg := NewPdbg(SetBuffers, OptSkips([]string{"customNo"}))
			customPdbgExcludeTest(apdbg)
			s := ShouldEqualNL(ErrString(), `r`)
			So(s, ShouldNotEqual, `e`)
		})
	})
}

func failCaller(skip int) (pc uintptr, file string, line int, ok bool) {
	return 0, "fail", skip, false
}
