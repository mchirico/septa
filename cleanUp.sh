#!/bin/bash
# Quick formatting tool

find . -iname '*.go' -exec gofmt -w {} \;
