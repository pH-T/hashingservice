package http

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

// NewHttpsServer returns a new HTTPS Server
func NewHttpsServer(hs hashservice, addr string, certPath, certKeyPath string) *httpsserver {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      initRoutes(hs),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	return &httpsserver{srv: srv, certPath: certPath, certKeyPath: certKeyPath}
}

type httpsserver struct {
	srv                   *http.Server
	mux                   http.Handler
	certPath, certKeyPath string
}

// Start starts the HTTP Server: blocks
func (server *httpsserver) Start() error {
	return server.srv.ListenAndServeTLS(server.certPath, server.certKeyPath)
}

// Stop stop the HTTP Server
func (server *httpsserver) Stop() error {
	log.Println("Closing...")
	return server.srv.Close()
}
