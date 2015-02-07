package exit

import "os"

// Func takes a code as exit status
type Func func(int)

// Exiter must have an Exit function able to get  an exit status
type Exiter interface {
	Exit(code int)
	Status() int
}

// Exit has an exiter, and will memorize the exit status code
type Exit struct {
	exit   Func
	status int
}

// Exit calls the exiter, and then returns code as status.
func (e *Exit) Exit(code int) {
	e.status = code
	e.exit(code)
}

// Status get the exit status code as memorized
// after the call to the exit func.
func (e *Exit) Status() int {
	return e.status
}

// DefaultExiter returns an exiter with default os.Exit() call.
// That means the status will never be visible,
// since os.Exit stops everything.
func DefaultExiter() Exiter {
	return &Exit{exit: os.Exit}
}

// NewExiter returns an exiter with a custom function
func NewExiter(exit Func) Exiter {
	return &Exit{exit: exit}
}
