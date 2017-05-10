#!/bin/sh
set -e

BINARY=goreleaser
FORMAT=tar.gz
OWNER=goreleaser
REPO=goreleaser
BINDIR=${BINDIR:-./bin}
TMPDIR=${TMPDIR:-/tmp}

VERSION=$1
if [ -z "${VERSION}" ]; then
  echo ""
  echo "Usage: $0 [version]"
  echo ""
  exit 1
fi
VERSION=${VERSION#v}

OS=$(uname -s)
ARCH=$(uname -m)



if [ ! -z "${ARM}" ]; then ARM="v$ARM"; fi
NAME=${BINARY}_${OS}_${ARCH}${ARM}
TARBALL=${NAME}.${FORMAT}
URL=https://github.com/${OWNER}/${REPO}/releases/download/v${VERSION}/${TARBALL}

if which curl > /dev/null; then
  WGET="curl -sSL -o"
elif which wget > /dev/null; then
  WGET="wget -q -O"
else
  echo "Unable to find wget or curl.  Exit"
  exit 1
fi

${WGET} ${TMPDIR}/${TARBALL} ${URL}
tar -C ${TMPDIR} -xzf ${TMPDIR}/${TARBALL}
install ${TMPDIR}/${BINARY} ${BINDIR}

