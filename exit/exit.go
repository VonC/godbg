package exit

// Inspired by
// http://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go

import "os"

// Func takes a code as exit status
type Func func(int)

// Exit has an exit func, and will memorize the exit status code
type Exit struct {
	exit   Func
	status int
}

// Exit calls the exiter, and then returns code as status.
// If e was declared, but never set (since only a test would set e),
// simply calls os.Exit()
func (e *Exit) Exit(code int) {
	if e != nil {
		e.status = code
		e.exit(code)
	} else {
		os.Exit(code)
	}
}

// Status get the exit status code as memorized
// after the call to the exit func.
func (e *Exit) Status() int {
	return e.status
}

// Default returns an Exit with default os.Exit() call.
// That means the status will never be visible,
// since os.Exit() stops everything.
func Default() *Exit {
	return &Exit{exit: os.Exit}
}

// NewExiter returns an exiter with a custom function
func New(exit Func) *Exit {
	return &Exit{exit: exit}
}
