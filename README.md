[![Build Status](https://travis-ci.org/mchirico/septa.svg?branch=develop)](https://travis-ci.org/mchirico/septa)
[![Go Report Card](https://goreportcard.com/badge/github.com/mchirico/septa)](https://goreportcard.com/report/github.com/mchirico/septa)

# SEPTA
Golang program for pulling SEPTA data

# Install
```bash


go get -u github.com/mchirico/septa/...

# To run it. This will create a token, if you don't have one.

septa

# This will just list one station at this time.

```

## Build
```bash
# Building locally

go build github.com/mchirico/septa/septa

```

## Go Packages for Doing Development
```bash
go get firebase.google.com/go

```


## Docker
This is early development; but, you can test it with a docker image

```bash
docker pull docker.io/mchirico/septa
docker run --rm -it mchirico/septa septa

```
