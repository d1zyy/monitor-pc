package pprofserver

import (
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
)

func StartPprofServer(addr string) error {
	mux := http.NewServeMux()
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting pprof server on %s\n", addr)

	return server.ListenAndServe()
}
