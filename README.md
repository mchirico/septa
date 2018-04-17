[![Build Status](https://travis-ci.org/mchirico/septa.svg?branch=develop)](https://travis-ci.org/mchirico/septa)
[![Maintainability](https://api.codeclimate.com/v1/badges/0282611068630ef5e232/maintainability)](https://codeclimate.com/github/mchirico/septa/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/mchirico/septa)](https://goreportcard.com/report/github.com/mchirico/septa)
[![Test Coverage](https://api.codeclimate.com/v1/badges/0282611068630ef5e232/test_coverage)](https://codeclimate.com/github/mchirico/septa/test_coverage)
[![codecov](https://codecov.io/gh/mchirico/septa/branch/develop/graph/badge.svg)](https://codecov.io/gh/mchirico/septa)

[![Waffle.io - Columns and their card count](https://badge.waffle.io/mchirico/septa.svg?columns=all)](https://waffle.io/mchirico/septa)
# SEPTA
Golang program for pulling SEPTA data

# Install
```bash

# install:

go get github.com/stretchr/testify/assert
go get firebase.google.com/go
go get -u github.com/mchirico/septa/...



```

## Build
```bash
# Building locally

go build github.com/mchirico/septa/septa

# To run it. This will create a token, if you don't have one.

septa

```

## Go Packages for Doing Development
```bash
go get firebase.google.com/go

```

## Server

This is the server program used to populate Firebase

```bash
go build github.com/mchirico/septa/routefirebase

```



## Docker
This is early development; but, you can test it with a docker image

```bash
docker pull docker.io/mchirico/septa
docker run --rm -it mchirico/septa septa

```
