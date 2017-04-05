package main

import (
	"flag"
	"github.com/ipcjk/ixgen/peeringdb"
	"log"
	"os"
)

func main() {
	var WaitForever chan struct{}
	signalChannels := make(chan os.Signal, 1)

	listenAPIServer := flag.String("listenAPI", "localhost:8443", "listenAddr for the api service")
	cacheDirectory := flag.String("cacheDir", "./cache", "cache directory for json files from peeringdb")
	_ = flag.String("style", "brocade/netiron", "Style for routing-config by template, e.g. brocade, juniper, cisco. Also possible: native/json or native/json_pretty for outputting the inside structures")
	_ = flag.String("templates", "./templates", "directory for templates")

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
