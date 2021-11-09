package net

import (
	"net/http"
	"net/http/pprof"
)

type pprofServer struct {
	Address string
	kind    string
}

func NewPProfServer(address string) *pprofServer {
	return &pprofServer{
		Address: address,
		kind:    PProfKind,
	}
}

func (p *pprofServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", indexHandler)
	mux.HandleFunc("/debug/pprof/heap", heapHandler)
	mux.HandleFunc("/debug/pprof/goroutine", goroutineHandler)
	mux.HandleFunc("/debug/pprof/allocs", allocsHandler)
	mux.HandleFunc("/debug/pprof/block", blockHandler)
	mux.HandleFunc("/debug/pprof/threadcreate", threadCreateHandler)
	mux.HandleFunc("/debug/pprof/cmdline", cmdlineHandler)
	mux.HandleFunc("/debug/pprof/profile", profileHandler)
	mux.HandleFunc("//debug/pprof/symbol", symbolHandler)
	mux.HandleFunc("/debug/pprof/trace", traceHandler)
	mux.HandleFunc("/debug/pprof/mutex", mutexHandler)
	err := http.ListenAndServe(p.Address, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *pprofServer) Stop() error {
	return nil
}

func (p *pprofServer) Kind() string {
	return p.kind
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

func heapHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("heap").ServeHTTP(w, r)
}
func goroutineHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("goroutine").ServeHTTP(w, r)
}
func allocsHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("allocs").ServeHTTP(w, r)
}
func blockHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("block").ServeHTTP(w, r)
}
func threadCreateHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("threadcreate").ServeHTTP(w, r)
}
func cmdlineHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Cmdline(w, r)
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Profile(w, r)
}
func symbolHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Symbol(w, r)
}
func traceHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Trace(w, r)
}
func mutexHandler(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("mutex").ServeHTTP(w, r)
}
