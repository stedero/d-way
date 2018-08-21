#! /bin/bash
echo "Go getting stuff"
go get github.com/rs/cors
# Alpine build
CGO_ENABLED=0 go build -a -installsuffix cgo d-way.go version.go
