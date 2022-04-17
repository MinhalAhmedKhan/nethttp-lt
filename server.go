package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {

	muxy := http.NewServeMux()
	muxy.Handle("/hello", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bye bye"))
		}),
	)

	var mux = NewMultiplexer()

	mux.Handle("/health", loggerMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("healthy"))
		}),
	)
	mux.Handle("/hello", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bye"))
		}),
	)

	mux.Handle("/bye", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bye"))
		}),
	)

	server := http.Server{
		Addr:              ":8080", // default is port 80
		Handler:           mux,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	fmt.Println("Server Starting 🚀")
	if err := server.ListenAndServe(); err != nil {
		fmt.Errorf("Failed to start server")
	}
}

// ---------- Handler ----------
type multiplexer struct {
	mu        sync.RWMutex // make endpoints thread safe
	endpoints map[string]http.Handler
}

func NewMultiplexer() *multiplexer {
	return &multiplexer{
		mu:        sync.RWMutex{},
		endpoints: make(map[string]http.Handler),
	}
}

func (mux *multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if handler, ok := mux.endpoints[r.URL.Path]; ok {
		handler.ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page not found"))
}

func (mux *multiplexer) Handle(uri string, handler http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.endpoints[uri] = handler
}

// ---------- Middlewares ----------
func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I got a request")
		next.ServeHTTP(w, r)
	}
}
