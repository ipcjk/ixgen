package main

import (
	"flag"
	"github.com/ipcjk/ixgen/libapiserver"
	"log"
	"os"
)

func main() {
	var Apiserver *libapiserver.Apiserver
	var WaitForever chan struct{}

	listenAPIServer := flag.String("listenAPI", "0.0.0.0:8443", "listenAddr for the api service")
	cacheDirectory := flag.String("cacheDir", "./cache", "cache directory for json files from peeringdb")
	templateDir := flag.String("templates", "./templates", "directory for templates")
	configDir := flag.String("configpath", "./configuration", "subdirectories for a user.tt to include general configuration commands, that will be expanded by the style - argument (not git tracked)")
	buildCache := flag.Bool("buildcache", false, "download json files for caching from peeringdb")

	flag.Parse()

	_, err := os.Stat(*cacheDirectory)
	if err != nil {
		log.Fatalf("Cant open cache directory: %s", err)
	}

	if *buildCache {
		libapiserver.DownloadCache("https://www.peeringdb.com/api", *cacheDirectory, "")
		os.Exit(0)
	}

	Apiserver = libapiserver.NewAPIServer(*listenAPIServer, *cacheDirectory, *templateDir, *configDir)
	Apiserver.RunAPIServer()
	<-WaitForever
}
