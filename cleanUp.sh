#!/bin/bash
# Quick formatting tool

find . -iname '*.go' -exec gofmt -w {} \;

# Tests
echo "Running tests..... "
go test github.com/mchirico/septa/firebase
go test github.com/mchirico/septa/utils

echo -e 'Want to test build?\n\n   go build github.com/mchirico/septa/routefirebase \n\n'

