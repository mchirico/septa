# SEPTA
Golang program for pulling SEPTA data

<a href="https://confluence.aipiggybot.io">
<img alt="Confluence" src="https://storage.googleapis.com/montco-stats/confluence.png"  width="100px">
</a>

<a href="https://jira.aipiggybot.io/projects/SEPT/issues/SEPT-5?filter=allopenissues">
<img alt="Confluence" src="https://storage.googleapis.com/montco-stats/jira.png"  width="22px">
</a>



[![Build Status](https://travis-ci.org/mchirico/septa.svg?branch=develop)](https://travis-ci.org/mchirico/septa)
[![Maintainability](https://api.codeclimate.com/v1/badges/0282611068630ef5e232/maintainability)](https://codeclimate.com/github/mchirico/septa/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/mchirico/septa)](https://goreportcard.com/report/github.com/mchirico/septa)
[![Test Coverage](https://api.codeclimate.com/v1/badges/0282611068630ef5e232/test_coverage)](https://codeclimate.com/github/mchirico/septa/test_coverage)
[![codecov](https://codecov.io/gh/mchirico/septa/branch/develop/graph/badge.svg)](https://codecov.io/gh/mchirico/septa)


Crossbrowser testing sponsored by [Browser Stack](https://www.browserstack.com)
[<img src="https://camo.githubusercontent.com/a7b268f2785656ab3ca7b1cbb1633ee5affceb8f/68747470733a2f2f64677a6f7139623561736a67312e636c6f756466726f6e742e6e65742f70726f64756374696f6e2f696d616765732f6c61796f75742f6c6f676f2d6865616465722e706e67" alt="Browser Stack" height="31px" style="background: cornflowerblue;">](https://www.browserstack.com)



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


## Map Command

```bash
# You may have a running version on p:
# Check this first
ssh p
tmux a -t 0

# No running version?
#  (Do the following from account "p"

tmux

docker pull us.gcr.io/mchirico/map

docker run -p 7444:80 --rm -it us.gcr.io/mchirico/map /bin/bash

/etc/init.d/postgresql start
/usr/sbin/apache2ctl start
/etc/init.d/renderd start



```



## Docker
This is early development; but, you can test it with a docker image

```bash
docker pull docker.io/mchirico/septa
docker run --rm -it mchirico/septa septa

```
