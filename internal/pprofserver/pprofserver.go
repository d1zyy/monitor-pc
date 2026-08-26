package pprofserver

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func StartPprofServer(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting pprof server on %s\n", addr)

	return server.ListenAndServe()
}
