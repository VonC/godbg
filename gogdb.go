package godbg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
)

var rxDbgLine, _ = regexp.Compile(`^.*[Vv]on[Cc](?:/prog/git)?/senvgo/main.go:(\d+)\s`)
var rxDbgFnct, _ = regexp.Compile(`^\s+(?:com/VonC/senvgo)?(?:\.\(([^\)]+)\))?\.?([^:]+)`)

// http://stackoverflow.com/a/23554672/6309 https://vividcortex.com/blog/2013/12/03/go-idiom-package-and-object/
// you design a type with methods as usual, and then you also place matching functions at the package level itself.
// These functions simply delegate to a default instance of the type that’s a private package-level variable, created in an init() function.

// Pdbg allows to print debug message with indent and function name added
type Pdbg struct {
	bout *bytes.Buffer
	berr *bytes.Buffer
	sout *bufio.Writer
	serr *bufio.Writer
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
var pdbg = &Pdbg{}

// Option set an option for a Pdbg
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*Pdbg)

// SetBuffers is an option for replacing stdout and stderr by
// bytes buffers (in a bufio.Writer)
func SetBuffers(apdbg *Pdbg) {
	apdbg.bout = bytes.NewBuffer(nil)
	apdbg.sout = bufio.NewWriter(apdbg.bout)
	apdbg.berr = bytes.NewBuffer(nil)
	apdbg.serr = bufio.NewWriter(apdbg.berr)
}

// NewPdbg creates a PDbg instance, with options
func NewPdbg(options ...Option) *Pdbg {
	newpdbg := &Pdbg{}
	for _, option := range options {
		option(newpdbg)
	}
	return newpdbg
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

// FlushIOs flushes the sout and serr bufio.Writer
func (pdbg *Pdbg) FlushIOs() {
	pdbg.sout.Flush()
	pdbg.serr.Flush()
}

func pdbgInc(scanner *bufio.Scanner, line string) string {
	m := rxDbgLine.FindSubmatchIndex([]byte(line))
	if len(m) == 0 {
		return ""
	}
	dbgLine := line[m[2]:m[3]]
	// fmt.Printf("line '%v', m '%+v'\n", line, m)
	scanner.Scan()
	line = scanner.Text()
	mf := rxDbgFnct.FindSubmatchIndex([]byte(line))
	// fmt.Printf("lineF '%v', mf '%+v'\n", line, mf)
	if len(mf) == 0 {
		return ""
	}
	dbgFnct := ""
	if mf[2] > -1 {
		dbgFnct = line[mf[2]:mf[3]]
	}
	if dbgFnct != "" {
		dbgFnct = dbgFnct + "."
	}
	dbgFnct = dbgFnct + line[mf[4]:mf[5]]

	return dbgFnct + ":" + dbgLine
}

func pdbgExcluded(dbg string) bool {
	if strings.Contains(dbg, "ReadConfig:") {
		return false
	}
	return false
}

// Pdbgf uses global Pdbg variable for printing strings, with indent and function name
func Pdbgf(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format+"\n", args...)
	msg = strings.TrimSpace(msg)
	bstack := bytes.NewBuffer(debug.Stack())
	// fmt.Printf("%+v", bstack)

	scanner := bufio.NewScanner(bstack)
	pmsg := ""
	depth := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "smartystreets") {
			break
		}
		m := rxDbgLine.FindSubmatchIndex([]byte(line))
		if len(m) == 0 {
			continue
		}
		if depth > 0 && depth < 4 {
			dbg := pdbgInc(scanner, line)
			if dbg == "" {
				continue
			}
			if depth == 1 {
				if pdbgExcluded(dbg) {
					return ""
				}
				pmsg = "[" + dbg + "]"
			} else {
				pmsg = pmsg + " (" + dbg + ")"
			}
		}
		depth = depth + 1
	}
	spaces := ""
	if depth >= 2 {
		spaces = strings.Repeat(" ", depth-2)
	}
	res := pmsg
	pmsg = spaces + pmsg
	msg = pmsg + "\n" + spaces + "  " + msg + "\n"
	// fmt.Printf("MSG '%v'\n", msg)
	fmt.Fprint(os.Stderr, fmt.Sprint(msg))
	return res
}
