package apiserverlib

import (
	"bytes"
	"github.com/ipcjk/ixgen/peeringdb"
	"io"
	"log"
	"net/http"
	"os"
)

var peeringDBfiles = []string{"fac", "ix", "ixfac", "ixlan", "ixpfx", "net", "netfac", "netixlan", "org", "poc"}

func DownloadCache(hostURL, cacheDir string) {
	for _, v := range peeringDBfiles {
		resp, err := http.Get(hostURL + "/" + v)
		if err != nil {
			log.Fatalf("Cant download %s", peeringDBfiles)
		}
		defer resp.Body.Close()

		targetFile := cacheDir + "/" + v
		writeCacheFile(targetFile+".download", resp.Body)
		data := readFile(targetFile + ".download")

		if v == "ix" {
			var apiResult peeringdb.Ix
			getJSON(bytes.NewBuffer(data), &apiResult)
			if len(apiResult.Data) < 500 {
				log.Fatalf("Cant update %s, missing records?", v)
			}
		} else if v == "net" {
			var apiResult peeringdb.Net
			getJSON(bytes.NewBuffer(data), &apiResult)
			if len(apiResult.Data) < 8000 {
				log.Fatalf("Cant update %s, missing records?", v)
			}
		} else if v == "netixlan" {
			var apiResult peeringdb.Netixlan
			getJSON(bytes.NewBuffer(data), &apiResult)
			if len(apiResult.Data) < 19000 {
				log.Fatalf("Cant update %s, missing records?", v)
			}
		}
		writeCacheFile(targetFile, bytes.NewBuffer(data))
		err = os.Remove(targetFile + ".download")
		if err != nil {
			log.Printf("Cant remove %s from fs", targetFile+".download")
		}
		log.Println("Updated " + targetFile)
		/*
			FIXME
			Check for meta-record and generation also
		*/

	}
}

func writeCacheFile(fileName string, reader io.Reader) {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalf("Cant open cache file target %s:", fileName)
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatalf("Could not copy file %s:", fileName)
	}
}
