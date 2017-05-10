package main

import (
	"text/template"
	"fmt"

	"github.com/goreleaser/goreleaser/config"
)

// converts the given name template to it's equivalent in shell
// except for the default goreleaser templates, templates with
// conditionals will return an error
//
// {{ .Binary }} --->  ${BINARY}, etc.
//
func nameTplInShell(target buildTarget) (string, error) {
	// TODO: error on conditionals
	if target == "{{ .Binary }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}" {
		prefix = "if [ ! -z \"${ARM}\" ]; then ARM=\"v$ARM\"; fi"
		target = "{{ .Binary }}_{{ .Os }}_{{ .Arch }}{{ .Arm }}"
	var varmap  = map[string]string{
		"Os": "${OS}",
		"Arch": "${ARCH}",
		"Arm": "${ARM}",
		"Version": "${VERSION}",
		"Tag": "${TAG}",
		"Binary": ${BINARY}",
	}
	
	var out bytes.Buffer
	t, err := template.New("name").Parse(target)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, varmap)
	return out.String(), err
}

var tplsrc = `#!/bin/sh
set -e

BINARY={{ .Binary }}

VERSION=$1
if [ -z "${VERSION}" ]; then
  echo ""
  echo "Usage: $0 [version]"
  echo ""
  exit 1
fi

OS=$(uname -s)
ARCH=$(uname -m)
VERSION=${VERSION#v}
BINDIR=${BINDIR:-./bin}
EXE=${BINDIR}/${BINARY}
TMPDIR=${TMPDIR:-/tmp}

case ${OS} in 

esac

case ${ARCH} in

esac

if [ ! -d "${BINDIR}" ]; then
  mkdir -p ${BINDIR}
fi

NAME={{ .Name }}

TARBALL=${NAME}.tar.gz
REPO=spf13/hugo
URL=https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}
echo "Downloading ${TARBALL}"
curl -sSL -o ${TMPDIR}/${TARBALL} ${URL}
tar -C ${TMPDIR} -xzf ${TMPDIR}/${TARBALL}
cp ${TMPDIR}/hugo ${HUGO}
`

func main() {

	

}
