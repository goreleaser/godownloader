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

shellcheck -s sh samples/*.sh
shellcheck -s bash samples/*.sh
shellcheck -s dash samples/*.sh
shellcheck -s ksh samples/*.sh
