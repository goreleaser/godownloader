#!/bin/sh
set -e

# Hugo binary downloader.
#
# Downloads the appropriate binary for macOS or linux
# to speed up your CI buids and to make sure your local Hugo
# is the same as the Hugo your CI/CD uses.
#
# Assumes curl is installed and `tar` supports the `-C` flag
# (change directory before executing)
#
# Why:
#
# 1. go get is slow
#
# "go get ...hugo..." can take up to 30 seconds on travis.ci.
# but downloading directly from GitHub releases takes 0.3 seconds.
#
# 2. making sure the local version is the same as travis.ci
#
# Best practice is to make sure your local development version of Hugo
# is the same as what your CI/CD system (e.g. travis, circle-ci, etc)
# is using.  This makes sure everyone is using the exact same version.
#
# FYI: URLs on GitHub look like this:
# https://github.com/spf13/hugo/releases/download/v0.20.7/hugo_0.20.7_Linux-64bit.tar.gz
#

# Set the version on the command line 
# as the first arg.
#
# ./download-hugo.sh spf13/hugo 0.20.0
#
PATH=$1
VERSION=$2
OWNER=$(dirname $PATH)
REPO=$(basename $PATH)
if [ -z "${OWNER}" or -z "${VERSION}"
  echo ""
  echo "Hugo binary downloader"
  echo ""
  echo "Usage:"
  echo "   $0: [owner/repo] [version]"
  echo "   where [owner/repo] is the GitHub.com identifier.
  echo "   where [version] is from the matching releases page.
  echo ""
  echo "   for example "
  echo "  $0 spf13/hugo 0.20.7  (see https://github.com/spf13/hugo/releases)"
  echo ""
  exit 1
fi


# strip leading "v" if any.
# tarball doesn't use "v", but URL does.
VERSION=${VERSION#v}

BINDIR="./bin"
mkdir -p ${BINDIR}
HUGO=${BINDIR}/hugo

# default TMPDIR to /tmp
# oddly not set on travis
TMPDIR=${TMPDIR:-/tmp}

# Do we need to get hugo?
#
#
if [ -f ${HUGO} ]; then
  if ${HUGO} version | grep -q ${VERSION}; then
     echo "Hugo ${VERSION} installed"
     exit 0
  fi
fi

# What OS?
#
# Translate what the OS claims it is to what
# Hugo GitHub releases uses
#
OS=$(uname -s)
case ${OS} in
Darwin)
  OS=macOS
  ;;
esac  

# Are we 64-bits or 32-bits?
#
#
CPU=$(uname -m)
case ${CPU} in
x86_64)
  CPU=64bit
  ;;
i386)
  CPU=32bit
  ;;
*) echo "unknown CPU: ${HOSTCPU}"
  exit 1
esac

TARBALL=hugo_${VERSION}_${OS}-${CPU}.tar.gz
REPO=spf13/hugo
URL=https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}
echo "Downloading ${TARBALL}"
curl -sSL -o ${TMPDIR}/${TARBALL} ${URL}
tar -C ${TMPDIR} -xzf ${TMPDIR}/${TARBALL}
cp ${TMPDIR}/hugo ${HUGO}
${HUGO} version

