#!/usr/bin/env bash
reset
cd ~/git/twine 
export GOPATH=~/git/twine 
go fmt gen
go build gen
go vet gen
./gen