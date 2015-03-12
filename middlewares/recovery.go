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

type recoveryBuilder func(*Recovery)

// Logger is a builder to define the logger used by the recover handler to log errors.
func Logger(logger *log.Logger) recoveryBuilder {
	return func(recovery *Recovery) {
		recovery.Logger = logger
	}
}

// PrintStackInBody is a builder to define if the recover handler should add the stack to the response body.
func PrintStackInBody(value bool) recoveryBuilder {
	return func(recovery *Recovery) {
		recovery.PrintStack = value
	}
}

// GetAllGoroutineStacks is a builder to define if the stack will contains only the current goroutine stack or all goroutines stacks
func GetAllGoroutineStacks(value bool) recoveryBuilder {
	return func(recovery *Recovery) {
		recovery.StackAllGoroutines = value
	}
}

// StackSize is a builder to set the size of the stack's buffer
func StackSize(value int) recoveryBuilder {
	return func(recovery *Recovery) {
		recovery.StackSize = value
	}
}

// DefaultRecovery instanciates the Recovery middleware with default values. Builders can update the default values.
func NewRecovery(builders ...recoveryBuilder) *Recovery {
	recovery := &Recovery{
		Logger:             log.New(os.Stdout, "", 0),
		PrintStack:         true,
		StackAllGoroutines: false,
		StackSize:          1024 * 8,
	}
	for _, builder := range builders {
		builder(recovery)
	}
	return recovery
}

// Recover is the MiddlewareBuilder function to use in the chain.
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
