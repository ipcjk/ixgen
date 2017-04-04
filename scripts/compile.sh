#!/bin/bash
echo "Building ixgen"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"  -o bin/ixgen.linux ixgen.go
env GOOS=windows GOARCH=amd64 go  build -ldflags="-s -w"  -o bin/ixgen.exe ixgen.go
env GOOS=darwin GOARCH=amd64 go  build -ldflags="-s -w"  -o bin/ixgen.mac ixgen.go
echo "Building apiserver"
env GOOS=linux GOARCH=amd64 go  build  -ldflags="-s -w" -o bin/apiserver.linux api_server/apiserver.go
env GOOS=windows GOARCH=amd64 go  build  -ldflags="-s -w" -o bin/apiserver.exe api_server/apiserver.go
env GOOS=darwin GOARCH=amd64 go  build  -ldflags="-s -w" -o bin/apiserver.mac api_server/apiserver.go
