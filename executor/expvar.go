// +build expvar

package executor

import (
	"expvar"
	"net"
	"net/http"
	"runtime"
)

// goroutines is an expvar.Func compliant wrapper for runtime.NumGoroutine function.
func goroutines() interface{} {
	return runtime.NumGoroutine()
}

func init() {
	expvar.Publish("Goroutines", expvar.Func(goroutines))

	stats = expvar.NewInt("GoroutinesInExecutor")

	lis, err := net.Listen(`tcp`, `localhost:5678`)
	if err != nil {
		panic(err)
	}

	up := make(chan *struct{})
	go func() {
		close(up)
		panic(http.Serve(lis, nil))
	}()
	<-up
}
