#! /bin/bash
# Alpine build
CGO_ENABLED=0 go build -a -installsuffix cgo d-way.go version.go