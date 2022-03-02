package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)


func main() {

	var mux = multiplexer{
		make(map[string]http.Handler),
	}
	mux.Handle("/health", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("healthy"))
		}),
	)
	mux.Handle("/hello", http.HandlerFunc(
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
		ConnContext:       modifyConnectionContext,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Errorf("Failed to start server")
	}

}


// ---------- Handler ----------
type multiplexer struct {
	endpoints map[string]http.Handler
}

func (m multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if handler, ok := m.endpoints[r.URL.Path]; ok {
		handler.ServeHTTP(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page not found"))
}

func (m multiplexer) Handle(uri string, handler http.Handler)  {
	m.endpoints[uri] = handler
}


// ---------- ConContext ----------
func modifyConnectionContext(ctx context.Context, c net.Conn) context.Context {
	return ctx
}