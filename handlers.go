package main

import "net/http"

const healthyResponse = "healthy"

func newHandlers() http.Handler{
	mux := http.NewServeMux()

	mux.Handle("/health", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(healthyResponse))
		}),
	)

	return mux
}
