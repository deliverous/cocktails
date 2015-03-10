package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

// Recovery is a middleware that recovers from any panics and writes a StatusInternalServerError.
type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// DefaultRecovery instanciates the Recovery middleware with default values.
func DefaultRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

// Recover is the MiddlewareBuilder function to use in the chain.
func (recovery *Recovery) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				stack := make([]byte, recovery.StackSize)
				stack = stack[:runtime.Stack(stack, recovery.StackAll)]

				f := "PANIC: %s\n%s"
				recovery.Logger.Printf(f, err, stack)

				if recovery.PrintStack {
					fmt.Fprintf(writer, f, err, stack)
				}
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
