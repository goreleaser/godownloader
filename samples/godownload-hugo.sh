#!/bin/sh
set -e

BINARY=hugo
FORMAT=tar.gz
OWNER=spf13
REPO=hugo
BINDIR=${BINDIR:-./bin}
TMPDIR=${TMPDIR:-/tmp}

VERSION=$1
if [ -z "${VERSION}" ]; then
echo "specify version number or 'latest'"
exit 1
fi

if [ "${VERSION}" = "latest" ]; then
  echo "Checking GitHub for latest version of ${OWNER}/${REPO}"
  VERSION=$(curl -s https://api.github.com/repos/${OWNER}/${REPO}/releases/latest | grep -m 1 "\"name\":" | cut -d ":" -f 2 | tr -d ' ",')
  if [ -z "${VERSION}" ]; then
    echo "Unable to determine latest release for ${OWNER}/${REPO}"
    exit 1
   fi
fi

VERSION=${VERSION#v}

OS=$(uname -s)
ARCH=$(uname -m)


case ${OS} in
Darwin) OS=macOS ;;
arm) OS=ARM ;;
arm64) OS=ARM64 ;;
dragonfly) OS=DragonFlyBSD ;;
i386) OS=32bit ;;
x86_64) OS=64bit ;;
esac

case ${ARCH} in
Darwin) ARCH=macOS ;;
arm) ARCH=ARM ;;
arm64) ARCH=ARM64 ;;
dragonfly) ARCH=DragonFlyBSD ;;
i386) ARCH=32bit ;;
x86_64) ARCH=64bit ;;
esac


NAME=${BINARY}_${VERSION}_${OS}-${ARCH}
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
install -d ${BINDIR}
install ${TMPDIR}/${BINARY} ${BINDIR}/

