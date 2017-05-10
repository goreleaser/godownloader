package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/goreleaser/goreleaser/config"
	yaml "gopkg.in/yaml.v1"
)

var tplsrc = `#!/bin/sh
set -e

BINARY={{ .Build.Binary }}
FORMAT={{ .Archive.Format }}
OWNER={{ $.Release.GitHub.Owner }}
REPO={{ $.Release.GitHub.Name }}
BINDIR=${BINDIR:-./bin}
TMPDIR=${TMPDIR:-/tmp}

VERSION=$1
if [ -z "${VERSION}" ]; then
echo "specify version number or 'latest'"
exit 1
fi

download() {
  DEST=$1
  SOURCE=$2
  if which curl > /dev/null; then
    WGET="curl -sSL"
    if [ "${DEST}" != "-" ]; then
      WGET="$WGET -o $DEST"
    fi
  elif which wget > /dev/null; then
    WGET="wget -q -O $DEST"
  else
    echo "Unable to find wget or curl.  Exit"
    exit 1
  fi

  # TODO: if source starts with github
  #  and we have env auth token
  #  then add it
  HEADER=""
  case $SOURCE in 
  https://api.github.com*)
     HEADER=""
     ;;
  esac
  ${WGET} $HEADER $SOURCE
}

if [ "${VERSION}" = "latest" ]; then
  echo "Checking GitHub for latest version of ${OWNER}/${REPO}"
  VERSION=$(download - https://api.github.com/repos/${OWNER}/${REPO}/releases/latest | grep -m 1 "\"name\":" | cut -d ":" -f 2 | tr -d ' ",')
  if [ -z "${VERSION}" ]; then
    echo "Unable to determine latest release for ${OWNER}/${REPO}"
    exit 1
   fi
fi

VERSION=${VERSION#v}

OS=$(uname -s)
ARCH=$(uname -m)

{{ with .Archive.Replacements }}
case ${OS} in 
{{- range $k, $v := . }}
{{ $k }}) OS={{ $v }} ;;
{{- end }}
esac

case ${ARCH} in
{{- range $k, $v := . }}
{{ $k }}) ARCH={{ $v }} ;;
{{- end }}
esac
{{ end }}

{{ .Archive.NameTemplate }}
TARBALL=${NAME}.${FORMAT}
URL=https://github.com/${OWNER}/${REPO}/releases/download/v${VERSION}/${TARBALL}

download ${TMPDIR}/${TARBALL} ${URL}
tar -C ${TMPDIR} -xzf ${TMPDIR}/${TARBALL}
install -d ${BINDIR}
install ${TMPDIR}/${BINARY} ${BINDIR}/
`

func makeShell(cfg *config.Project) (string, error) {
	var out bytes.Buffer
	t, err := template.New("shell").Parse(tplsrc)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, cfg)
	return out.String(), err
}

// converts the given name template to it's equivalent in shell
// except for the default goreleaser templates, templates with
// conditionals will return an error
//
// {{ .Binary }} --->  NAME=${BINARY}, etc.
//
func makeName(target string) (string, error) {
	prefix := ""
	// TODO: error on conditionals
	if target == "" || target == "{{ .Binary }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}" {
		prefix = "if [ ! -z \"${ARM}\" ]; then ARM=\"v$ARM\"; fi"
		target = "{{ .Binary }}_{{ .Os }}_{{ .Arch }}{{ .Arm }}"
	}
	var varmap = map[string]string{
		"Os":      "${OS}",
		"Arch":    "${ARCH}",
		"Arm":     "${ARM}",
		"Version": "${VERSION}",
		"Tag":     "${TAG}",
		"Binary":  "${BINARY}",
	}

	var out bytes.Buffer
	if prefix != "" {
		out.WriteString(prefix + "\n")
	}
	out.WriteString("NAME=")
	t, err := template.New("name").Parse(target)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, varmap)
	return out.String(), err
}

func readURL(loc string) ([]byte, error) {
	resp, err := http.Get(loc)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	// to make errcheck be happy
	errc := resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if errc != nil {
		return nil, errc
	}
	return body, err
}

func Load(repo string, file string) (*config.Project, error) {
	if repo == "" && file == "" {
		return nil, fmt.Errorf("Need a repo or file")
	}
	if file == "" {
		file = "https://raw.githubusercontent.com/" + repo + "/master/goreleaser.yml"
	}
	var body []byte
	var err error
	log.Printf("Reading %s", file)
	if strings.HasPrefix(file, "http") {
		body, err = readURL(file)
	} else {
		body, err = ioutil.ReadFile(file)
	}
	if err != nil {
		return nil, err
	}
	project := &config.Project{}
	err = yaml.Unmarshal(body, project)
	if err != nil {
		return nil, err
	}

	// if not specified add in GitHub owner/repo info
	if project.Release.GitHub.Owner == "" {
		if repo == "" {
			return nil, fmt.Errorf("Need to provide owner/name repo!")
		}
		project.Release.GitHub.Owner = path.Dir(repo)
		project.Release.GitHub.Name = path.Base(repo)
	}

	// set default archive format
	if project.Archive.Format == "" {
		project.Archive.Format = "tar.gz"
	}

	// set default binary name
	if project.Build.Binary == "" {
		project.Build.Binary = path.Base(repo)
	}

	// Convert replacements from GOOS/GOARCH to uname.

	// map of golang OS/ARCH identifier to what uname uses
	uname := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"freebsd": "FreeBSD",
		"openbsd": "OpenBSD",
		"netbsd":  "NetBSD",
		"windows": "Windows",
		"386":     "i386",
		"amd64":   "x86_64",
	}
	rmap := make(map[string]string)
	for k, v := range project.Archive.Replacements {
		newk := uname[k]

		// if unknown, keep
		if newk == "" {
			rmap[k] = v
			continue
		}

		// if mapping is an idenity, then ignore
		if newk == v {
			continue
		}

		rmap[newk] = v
	}
	project.Archive.Replacements = rmap

	return project, nil
}

func main() {
	repo := flag.String("repo", "", "owner/name of repository")
	flag.Parse()
	args := flag.Args()
	file := ""
	if len(args) > 0 {
		file = args[0]
	}
	cfg, err := Load(*repo, file)
	if err != nil {
		log.Fatalf("Unable to parse: %s", err)
	}

	// get name template
	name, err := makeName(cfg.Archive.NameTemplate)
	cfg.Archive.NameTemplate = name
	if err != nil {
		log.Fatalf("Unable generate name: %s", err)
	}

	shell, err := makeShell(cfg)
	if err != nil {
		log.Fatalf("Unable to generate shell: %s", err)
	}
	fmt.Println(shell)
}
