package http

import (
	"context"
	"log"
	"net/http"
	"time"
)

type hashservice interface {
	Hash(context.Context, string) (string, error)
	Verify(context.Context, string, string) (bool, bool, error)
}

// NewHttpServer returns a new HTTP Server
func NewHttpServer(hs hashservice, addr string) *httpserver {
	srv := &http.Server{
		Addr:         addr,
		Handler:      initRoutes(hs),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &httpserver{srv: srv}
}

type httpserver struct {
	srv *http.Server
	mux http.Handler
}

// Start starts the HTTP Server: blocks
func (server *httpserver) Start() error {
	return server.srv.ListenAndServe()
}

// Stop stop the HTTP Server
func (server *httpserver) Stop() error {
	log.Println("Closing...")
	return server.srv.Close()
}
