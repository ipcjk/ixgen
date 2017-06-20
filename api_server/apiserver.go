package main

import (
	"flag"
	"github.com/ipcjk/ixgen/libapiserver"
	"log"
	"os"
)

func main() {
	var WaitForever chan struct{}
	signalChannels := make(chan os.Signal, 1)

	listenAPIServer := flag.String("listenAPI", "0.0.0.0:8443", "listenAddr for the api service")
	cacheDirectory := flag.String("cacheDir", "./cache", "cache directory for json files from peeringdb")
	templateDir := flag.String("templates", "./templates", "directory for templates")

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

	Apiserver := libapiserver.NewAPIServer(*listenAPIServer, *cacheDirectory, *templateDir)
	Apiserver.RunAPIServer()
	<-WaitForever
}
