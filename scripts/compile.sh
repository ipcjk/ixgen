#!/bin/bash
echo "Building ixgen"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -o ixgen.linux ixgen.go
env GOOS=windows GOARCH=amd64 go  build -ldflags="-s -w"  -o ixgen.exe ixgen.go
env GOOS=darwin GOARCH=amd64 go  build -ldflags="-s -w"  -o ixgen.mac ixgen.go
echo "Building apiserver"
env GOOS=linux GOARCH=amd64 go  build  -ldflags="-s -w" -o apiserver.linux api_server/apiserver.go
env GOOS=windows GOARCH=amd64 go  build  -ldflags="-s -w" -o apiserver.exe api_server/apiserver.go
env GOOS=darwin GOARCH=amd64 go  build  -ldflags="-s -w" -o apiserver.mac api_server/apiserver.go
