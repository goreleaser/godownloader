#!/bin/sh -ex

gometalinter --vendor --disable-all \
  --enable=deadcode \
  --enable=ineffassign \
  --enable=gosimple \
  --enable=staticcheck \
  --enable=gofmt \
  --enable=goimports \
  --enable=dupl \
  --enable=misspell \
  --enable=errcheck \
  --enable=vetshadow \
  --deadline=10m \
  ./...

# shellcheck put into seperate file since travis-ci is broken
# https://github.com/goreleaser/godownloader/issues/61
