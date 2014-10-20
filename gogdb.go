package godbg

import (
	"bufio"
	"bytes"
	"fmt"
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
}

// global pdbg used for printing
var pdbg *Pdbg = &Pdbg{}

// Option set an option for a Pdbg
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type Option func(*Pdbg) error

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
