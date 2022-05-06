package mq

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	ErrQueueClosed         = errors.New("mq: Queue is closed")
	ErrQueueEmtpy          = errors.New("mq: Queue is empty")
	ErrQueueFilled         = errors.New("mq: Queue is full")
	ErrQueueUnavailable    = errors.New("mq: Queue not found")
	ErrConsumerUnavailable = errors.New("mq: Consumer not found")
)

// This output is shown if a panic happens.
const panicOutput = `
!!!!!!!!!!!!!!!!!!!!!!!!!!! QUEEING SERVICE CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
Queuing Service crashed! This is always indicative of a bug within the service.
!!!!!!!!!!!!!!!!!!!!!!!!!!! QUEEING SERVICE CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
`

// In case multiple goroutines panic concurrently, ensure only the first one
// recovered by PanicHandler starts printing.
var panicMutex sync.Mutex

// PanicHandler is called to recover from an internal panic within service, and
// augments the standard stack trace with a more user friendly error message.
// PanicHandler must be called as a defered function, and must be the first
// defer called at the start of a new goroutine.
func PanicHandler() {
	// Have all managed goroutines checkin here, and prevent them from exiting
	// if there's a panic in progress. While this can't lock the entire runtime
	// to block progress, we can prevent some cases where Terraform may return
	// early before the panic has been printed out.
	panicMutex.Lock()
	defer panicMutex.Unlock()

	recovered := recover()
	if recovered == nil {
		return
	}

	fmt.Fprint(os.Stderr, panicOutput)
	fmt.Fprint(os.Stderr, recovered, "\n")

	// An exit code of 11 keeps us out of the way of the detailed exitcodes
	// from plan, and also happens to be the same code as SIGSEGV which is
	// roughly the same type of condition that causes most panics.
	os.Exit(11)
}
