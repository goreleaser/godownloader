#!/bin/sh
set -ex
go get -u github.com/alecthomas/gometalinter
go get -u github.com/golang/dep/...
go get -u github.com/pierrre/gotestcover
go get -u golang.org/x/tools/cmd/cover
dep ensure
gometalinter --install
