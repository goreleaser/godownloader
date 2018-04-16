#!/bin/sh -ex

gometalinter.v2 --vendor ./...

# commented because shellcheck was put into seperate file since travis-ci is broken
# https://github.com/goreleaser/godownloader/issues/61
