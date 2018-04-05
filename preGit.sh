#!/bin/bash
# Quick formatting tool

find . -iname '*.go' -exec gofmt -w {} \;

# Tests
echo "Running tests..... "
go test -coverprofile=c.out -covermode=atomic github.com/mchirico/septa/utils	github.com/mchirico/septa/firebase
go vet github.com/mchirico/septa/utils	github.com/mchirico/septa/firebase

echo -e 'Want to test build?\n\n   go build github.com/mchirico/septa/routefirebase \n\n'

