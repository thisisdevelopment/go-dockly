package xmiddleware

import (
	"fmt"
	"github.com/bugsnag/bugsnag-go/v2"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
)

// PanicReceiver is able to receive panic information from the Recoverer middleware
type PanicReceiver interface {
	ReceivePanic(msg any, stack []byte)
}

// PanicReceiverFunction function type that implements PanicReceiver interface
type PanicReceiverFunction func(any, []byte)

// ReceivePanic receives panic information whenever the Recoverer middleware catches a panic
func (f PanicReceiverFunction) ReceivePanic(msg any, stack []byte) {
	f(msg, stack)
}

// BugsnagPanicReceiver receives panic information and notifies bugsnag
type BugsnagPanicReceiver struct{}

// NewBugsnagPanicReceiver creates a new bugsnag panic receiver
func NewBugsnagPanicReceiver() *BugsnagPanicReceiver {
	return &BugsnagPanicReceiver{}
}

// ReceivePanic receives panic information whenever the Recoverer middleware catches a panic
func (bs *BugsnagPanicReceiver) ReceivePanic(msg any, _ []byte) {
	_ = bugsnag.Notify(fmt.Errorf("%v", msg), bugsnag.SeverityError, bugsnag.Context{String: "panic"})
}

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func Recoverer(logLevel string, receivers ...PanicReceiver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					stackTrace := debug.Stack()

					// pass panic and stack to receivers
					for _, receiver := range receivers {
						receiver.ReceivePanic(rvr, stackTrace)
					}

					logEntry := middleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rvr, stackTrace)
					} else {
						_, _ = fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
						debug.PrintStack()
					}

					errorText := http.StatusText(http.StatusInternalServerError)

					if logLevel == "debug" {
						errorText = fmt.Sprintf("%s\n %+v\n", errorText, rvr)
					}

					http.Error(w, errorText, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
