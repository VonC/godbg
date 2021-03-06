package godbg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// http://stackoverflow.com/a/23554672/6309 https://vividcortex.com/blog/2013/12/03/go-idiom-package-and-object/
// you design a type with methods as usual, and then you also place matching functions at the package level itself.
// These functions simply delegate to a default instance of the type that’s a private package-level variable, created in an init() function.

// Pdbg allows to print debug message with indent and function name added
type Pdbg struct {
	bout     *bytes.Buffer
	berr     *bytes.Buffer
	sout     *bufio.Writer
	serr     *bufio.Writer
	breaks   []string
	excludes []string
	skips    []string
}

// Out returns a writer for normal messages.
// By default, os.StdOut
func Out() io.Writer {
	return pdbg.Out()
}

// Out returns a writer for normal messages for a given pdbg instance.
// By default, os.StdOut
func (pdbg *Pdbg) Out() io.Writer {
	if pdbg.sout == nil {
		return os.Stdout
	}
	return pdbg.sout
}

// Err returns a writer for error messages.
// By default, os.StdErr
func Err() io.Writer {
	return pdbg.Err()
}

// Err returns a writer for error messages for a given pdbg instance.
// By default, os.StdErr
func (pdbg *Pdbg) Err() io.Writer {
	if pdbg.serr == nil {
		return os.Stderr
	}
	return pdbg.serr
}

// global pdbg used for printing
var pdbg = NewPdbg()

// Option set an option for a Pdbg
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*Pdbg)

// SetBuffers is an option for replacing stdout and stderr by
// bytes buffers (in a bufio.Writer).
// If apdbg is nil, set for the global pdbg instance
func SetBuffers(apdbg *Pdbg) {
	if apdbg == nil {
		apdbg = pdbg
	}
	apdbg.bout = bytes.NewBuffer(nil)
	apdbg.sout = bufio.NewWriter(apdbg.bout)
	apdbg.berr = bytes.NewBuffer(nil)
	apdbg.serr = bufio.NewWriter(apdbg.berr)
}

// SetExcludes set excludes on a pdbg (nil for global pdbg)
func (pdbg *Pdbg) SetExcludes(excludes []string) {
	pdbg.excludes = excludes
}

// OptExcludes is an option to set excludes at the creation of a pdbg
func OptExcludes(excludes []string) Option {
	return func(apdbg *Pdbg) {
		apdbg.SetExcludes(excludes)
	}
}

// SetSkips set skips on a pdbg (nil for global pdbg)
func (pdbg *Pdbg) SetSkips(skips []string) {
	sk := []string{"/godbg.go'"}
	sk = append(sk, skips...)
	pdbg.skips = sk
}

// OptSkips is an option to set excludes at the creation of a pdbg
func OptSkips(skips []string) Option {
	return func(apdbg *Pdbg) {
		apdbg.SetSkips(skips)
	}
}

// NewPdbg creates a PDbg instance, with options
func NewPdbg(options ...Option) *Pdbg {
	newpdbg := &Pdbg{}
	for _, option := range options {
		option(newpdbg)
	}
	newpdbg.breaks = append(newpdbg.breaks, "smartystreets")
	//newpdbg.breaks = append(newpdbg.breaks, "(*Pdbg).Pdbgf")
	newpdbg.skips = append(newpdbg.skips, "/godbg.go'")
	return newpdbg
}

// ResetIOs reset the out and err buffer of global pdbg instance
func ResetIOs() {
	pdbg.ResetIOs()
}

// ResetIOs reset the out and err buffer
// (unless they were the default stdout and stderr,
// in which case it does nothing)
func (pdbg *Pdbg) ResetIOs() {
	if pdbg.sout != nil {
		pdbg.bout = bytes.NewBuffer(nil)
		pdbg.sout.Reset(pdbg.bout)
		pdbg.berr = bytes.NewBuffer(nil)
		pdbg.serr.Reset(pdbg.berr)
	}
}

// OutString returns the string for out messages for the global pdbg instance.
// It flushes the out buffer.
// If out is set to os.Stdout, returns an empty string
func OutString() string {
	return pdbg.OutString()
}

// OutString returns the string for out messages for a given pdbg instance.
// It flushes the out buffer.
// If out is set to os.Stdout, returns an empty string
func (pdbg *Pdbg) OutString() string {
	if pdbg.sout == nil {
		return ""
	}
	pdbg.sout.Flush()
	return pdbg.bout.String()
}

// ErrString returns the string for error messages for the global pdbg instance.
// It flushes the err buffer.
// If err is set to os.StdErr, returns an empty string
func ErrString() string {
	return pdbg.ErrString()
}

// ErrString returns the string for error messages for a given pdbg instance.
// It flushes the err buffer.
// If err is set to os.StdErr, returns an empty string
func (pdbg *Pdbg) ErrString() string {
	if pdbg.serr == nil {
		return ""
	}
	pdbg.serr.Flush()
	return pdbg.berr.String()
}

// NoOutput checks if there is any output recorded on Stdout or Stderr
func NoOutput() bool {
	return OutString() == "" && ErrString() == ""
}

// NoOutput checks if there is any output recorded on Stdout or Stderr
// for a given pdbg instance.
func (pdbg *Pdbg) NoOutput() bool {
	return pdbg.OutString() == "" && pdbg.ErrString() == ""
}

func (pdbg *Pdbg) pdbgExcluded(dbg string) bool {
	for _, e := range pdbg.excludes {
		if strings.Contains(dbg, e) {
			// fmt.Printf("EXCLUDE over '%v' including '%v'\n", dbg, e) // DBG
			return true
		}
	}
	return false
}

func (pdbg *Pdbg) pdbgBreak(dbg string) bool {
	for _, b := range pdbg.breaks {
		if strings.Contains(dbg, b) {
			// fmt.Printf("BREAK over '%v' including '%v'\n", dbg, b) // DBG
			return true
		}
	}
	return false
}

func (pdbg *Pdbg) pdbgSkip(dbg string) (bool, int) {
	depthToAdd := 0
	for i, s := range pdbg.skips {
		if strings.Contains(dbg, s) {
			if i > 0 {
				depthToAdd = 1
			}
			// fmt.Printf("SKIP over '%v' including '%v'\n", dbg, s) // DBG
			return true, depthToAdd
		}
	}
	return false, depthToAdd
}

// Pdbgf uses global Pdbg variable for printing strings, with indent and function name
func Pdbgf(format string, args ...interface{}) string {
	return pdbg.Pdbgf(format, args...)
}

type caller func(skip int) (pc uintptr, file string, line int, ok bool)

var mycaller = runtime.Caller

// Perrdbgf uses Stderr for printing strings, with indent and function name
func Perrdbgf(format string, args ...interface{}) string {
	return pdbg.pdbgfw(format, os.Stderr, args...)
}

// Pdbgf uses custom Pdbg variable for printing strings, with indent and function name
func (pdbg *Pdbg) Pdbgf(format string, args ...interface{}) string {
	return pdbg.pdbgfw(format, pdbg.Err(), args...)
}

func (pdbg *Pdbg) pdbgfw(format string, iow io.Writer, args ...interface{}) string {
	msg := fmt.Sprintf(format+"\n", args...)
	msg = strings.TrimSpace(msg)

	pmsg := ""
	depth := 0
	nbskip := 0
	nbInitialSkips := 0
	first := true
	addOneForSkip := 0
	// fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~\n") // DBG
	for ok := true; ok; {
		pc, file, line, ok := mycaller(depth)
		if !ok {
			break
		}
		fname := runtime.FuncForPC(pc).Name()
		fline := fmt.Sprintf("Name of function: '%v': '%+x' (line %v): file '%v'\n", fname, fname, line, file)
		// fmt.Println(fline) // DBG
		if pdbg.pdbgExcluded(fline) {
			depth = depth + 1
			if first {
				return ""
			}
			continue
		}
		if pdbg.pdbgBreak(fline) {
			break
		}
		if isSkipped, depthToAdd := pdbg.pdbgSkip(fline); isSkipped {
			depth = depth + 1
			nbskip = nbskip + 1
			addOneForSkip = addOneForSkip + depthToAdd
			continue
		}
		fnamerx1 := regexp.MustCompile(`.*\.func[^a-zA-Z0-9]`)
		fname = fnamerx1.ReplaceAllString(fname, "func.")
		fnamerx2 := regexp.MustCompile(`.*/`)
		fname = fnamerx2.ReplaceAllString(fname, "")
		if !strings.HasPrefix(fname, "func.") {
			fnamerx3 := regexp.MustCompile(`^.*?\.`)
			// fmt.Printf("fname before: '%v'", fname)
			fname = fnamerx3.ReplaceAllString(fname, "")
			// fmt.Printf(" => fname after: '%v'\n", fname)
			fnamerx4 := regexp.MustCompile(`[\(\)]`)
			fname = fnamerx4.ReplaceAllString(fname, "")
		}
		dbg := fname + ":" + fmt.Sprintf("%d", line)
		if first {
			// fmt.Printf(" => nbskip '%v'; addOneForSkip '%v'\n", nbskip, addOneForSkip) // DBG
			nbInitialSkips = nbskip - addOneForSkip
			pmsg = "[" + dbg + "]"
		} else {
			pmsg = pmsg + " (" + dbg + ")"
		}
		first = false
		depth = depth + 1
	}
	finalDepth := depth
	depth = finalDepth - nbInitialSkips - 1

	spaces := ""
	if depth >= 0 {
		spaces = strings.Repeat(" ", depth*2)
	}
	// fmt.Printf("spaces '%s', finalDepth '%d', depth '%d', nbInitialSkips '%d', addOneForSkip='%d'\n", spaces, finalDepth, depth, nbInitialSkips, addOneForSkip) // DBG
	res := pmsg
	if pmsg != "" {
		pmsg = spaces + pmsg + "\n"
	}
	msg = pmsg + spaces + "  " + msg + "\n"
	// fmt.Printf("==> MSG '%v'\n", msg) // DBG
	fmt.Fprint(iow, fmt.Sprint(msg))
	return res
}

var r = regexp.MustCompile(`:\d+[\)\]]`)
var r2 = regexp.MustCompile(`func\.\d+[\)\]]`)

// ShouldEqualNL is a custom goconvey assertion to ignore differences
// with func id and lines: `[globalPdbgExcludeTest:16] (func.019:167)` would
// be equal to [globalPdbgExcludeTest] (func)
// (see https://github.com/smartystreets/goconvey/wiki/Custom-Assertions)
func ShouldEqualNL(actual interface{}, expected ...interface{}) string {
	a := actual.(string)
	e := expected[0].(string)
	a = r.ReplaceAllStringFunc(a, func(s string) string { return s[len(s)-1:] })
	e = r.ReplaceAllStringFunc(e, func(s string) string { return s[len(s)-1:] })
	a = r2.ReplaceAllStringFunc(a, func(s string) string { return "func" + s[len(s)-1:] })
	e = r2.ReplaceAllStringFunc(e, func(s string) string { return "func" + s[len(s)-1:] })
	a = strings.TrimRight(a, "\r\n")
	e = strings.TrimRight(e, "\r\n")
	a = strings.Replace(a, "\r\n", "\n", -1)
	if a == e {
		return ""
	}
	al := strings.Split(a, "\n")
	el := strings.Split(e, "\n")
	msg := ""
	diff := false
	ali := 0
	aline := ""
	for ali, aline = range al {
		if len(el) <= ali {
			msg = msg + fmt.Sprintf("Expected has only %d lines instead of actual %d lines: %d lines not expected", ali, len(al), len(al)-ali)
			diff = true
			break
		}
		eline := el[ali]
		if aline == eline {
			continue
		}
		msg = msg + fmt.Sprintf("line '%d' differs:\nE'%s'\nA'%s'\nE'%q'\nA'%q'", ali+1, eline, aline, eline, aline)
		diff = true
		break
	}
	if !diff && len(el) >= ali+1 {
		msg = msg + fmt.Sprintf("actual misses lines after line %d\n'%s'", ali+1, el[ali+1])
	}
	return fmt.Sprintf("Expected: '%s'\nActual:   '%s'\n(Should be equal even with different lines and function ids)\n%s", e, a, msg)
}
