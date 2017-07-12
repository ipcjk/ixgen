#!/bin/bash
echo "Building ixgen"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -o bin/ixgen.linux ixgen.go
env GOOS=windows GOARCH=amd64 go  build -ldflags="-s -w"  -o bin/ixgen.exe ixgen.go
env GOOS=darwin GOARCH=amd64 go  build -ldflags="-s -w"  -o bin/ixgen.mac ixgen.go
echo "Building ixapiserver"
env GOOS=linux GOARCH=amd64 go  build  -ldflags="-s -w" -o bin/ixapiserver.linux ixapiserver/ixapiserver.go
env GOOS=windows GOARCH=amd64 go  build  -ldflags="-s -w" -o bin/ixapiserver.exe ixapiserver/ixapiserver.go
env GOOS=darwin GOARCH=amd64 go  build  -ldflags="-s -w" -o bin/ixapiserver.mac ixapiserver/ixapiserver.go
