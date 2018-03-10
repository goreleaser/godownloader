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

# SC2034 is unused variable 
# some generated scripts contain 1 or more variables with aren't used
# sometimes.
shellcheck -e SC2034 -s sh samples/*.sh
shellcheck -e SC2034 -s bash samples/*.sh
shellcheck -e SC2034 -s dash samples/*.sh
shellcheck -e SC2034 -s ksh samples/*.sh
