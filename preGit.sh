#!/bin/bash
# Quick formatting tool

find . -iname '*.go' -exec gofmt -w {} \;

# Tests
echo "Running tests..... "
go test -coverprofile=coverage.txt -covermode=atomic github.com/mchirico/septa/firebase
go test -coverprofile=coverage.txt -covermode=atomic github.com/mchirico/septa/utils
go test -cover fmt
echo -e 'Want to test build?\n\n   go build github.com/mchirico/septa/routefirebase \n\n'

