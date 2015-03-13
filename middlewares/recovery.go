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
	Logger             *log.Logger
	PrintStack         bool
	StackAllGoroutines bool
	StackSize          int
}

// SetLogger defines the logger used by the recover handler to log errors.
func (recovery *Recovery) SetLogger(logger *log.Logger) *Recovery {
	recovery.Logger = logger
	return recovery
}

// SetPrintStackInBody defines if the recover handler should add the stack to the response body.
func (recovery *Recovery) SetPrintStackInBody(value bool) *Recovery {
	recovery.PrintStack = value
	return recovery
}

// GetAllGoroutineStacks is a builder to define if the stack will contains only the current goroutine stack or all goroutines stacks
func (recovery *Recovery) GetAllGoroutineStacks(value bool) *Recovery {
	recovery.StackAllGoroutines = value
	return recovery
}

// SetStackSize builders to set the size of the stack's buffer
func (recovery *Recovery) SetStackSize(value int) *Recovery {
	recovery.StackSize = value
	return recovery
}

// DefaultRecovery instanciates the Recovery middleware with default values.
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:             log.New(os.Stdout, "", 0),
		PrintStack:         true,
		StackAllGoroutines: false,
		StackSize:          1024 * 8,
	}
}

// Recover is the Middleware function to use in the chain.
func (recovery *Recovery) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				stack := make([]byte, recovery.StackSize)
				stack = stack[:runtime.Stack(stack, recovery.StackAllGoroutines)]

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
