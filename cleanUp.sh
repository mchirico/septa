#!/bin/bash
# Quick formatting tool

find . -iname '*.go' -exec gofmt -w {} \;

# Tests
echo "Running tests..... "
go test github.com/mchirico/septa/firebase
go test github.com/mchirico/septa/utils

