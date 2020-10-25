package main

import (
	"flag"
	"fmt"
	"hashservice/server/http"
	"hashservice/service"
	"hashservice/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var version string = "hashservice v1.0.0"
var addr string = ""
var debug bool = false
var certPath, certKeyPath string

func init() {
	flag.StringVar(&addr, "l", "localhost:8080", "Where the server should listen")
	flag.StringVar(&certPath, "c", "", "Server certificate")
	flag.StringVar(&certKeyPath, "ck", "", "Server certificate key")
	flag.BoolVar(&debug, "d", false, "Debug mode: no HTTPS")
	flag.Parse()
}

type Server interface {
	Start() error
	Stop() error
}

func main() {
	fmt.Printf("Running: %s\n", version)

	store := storage.NewMemoryStorage()
	hashservice := service.NewBcryptHashService(store)

	var server Server
	if !debug {
		if certPath == "" || certKeyPath == "" {
			log.Fatal("Certificat-Info for HTTPS is missing!")
			flag.Usage()
			os.Exit(1)
		}
		server = http.NewHttpsServer(hashservice, addr, certPath, certKeyPath)
	} else {
		log.Println("DEBUG")
		server = http.NewHttpServer(hashservice, addr)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Stop()
	}()

	log.Printf("Starting on: %s", addr)
	log.Fatal(server.Start())

}
