#!/usr/bin/env bash
GOOS=windows GOARCH=386 go build -o bin/placeimag_386.exe placeimag.go
GOOS=windows GOARCH=amd64 go build -o bin/placeimag.exe placeimag.go
GOOS=darwin GOARCH=amd64 go build -o bin/placeimag_mac placeimag.go
go build -o bin/placeimag placeimag.go
