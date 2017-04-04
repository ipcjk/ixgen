package main

import (
	"IXgenerator/peeringdb"
	"flag"
	"log"
	"os"
)

func main() {
	var WaitForever chan struct{}
	signalChannels := make(chan os.Signal, 1)

	listenAPIServer := flag.String("listenAPI", "localhost:8443", "listenAddr for the api service")
	cacheDirectory := flag.String("cacheDir", "./cache", "cache directory for json files from peeringdb")
	flag.Parse()

	go func() {
		for range signalChannels {
			// e.g. throw away cache?
			// currently: exit to os
			WaitForever <- struct{}{}
		}
	}()

	_, err := os.Stat(*cacheDirectory)
	if err != nil {
		log.Fatalf("Cant open cache directory: %s", err)
	}

	Apiserver := peeringdb.NewAPIServer(*listenAPIServer, *cacheDirectory)
	Apiserver.RunAPIServer()
	<-WaitForever
}
