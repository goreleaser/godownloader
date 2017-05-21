package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/goreleaser/goreleaser/config"
	"github.com/goreleaser/goreleaser/pipeline/defaults"
)

// given a template, and a config, generate shell script
func makeShell(tplsrc string, cfg *config.Project) (string, error) {
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
	if target == "" {
		target = defaults.NameTemplate
	}

	// armv6 is the default in the shell script
	// so do not need special template condition for ARM
	armversion := "{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}"
	target = strings.Replace(target, armversion, "{{ .Arch }}", -1)

	// otherwise if it contains a conditional, we can't (easily)
	// translate that to bash.  Ask for bug report.
	if strings.Contains(target, "{{ if") || strings.Contains(target, "{{if") || strings.Contains(target, "{{ .Arm") || strings.Contains(target, "{{.Arm") {
		return "", fmt.Errorf("name_template %q contains unknown conditional or ARM format.  Please file bug at https://github.com/goreleaser/godownloader", target)
	}

	var varmap = map[string]string{
		"Os":      "${OS}",
		"Arch":    "${ARCH}",
		"Version": "${VERSION}",
		"Tag":     "${TAG}",
		"Binary":  "${BINARY}",
	}

	var out bytes.Buffer
	out.WriteString("NAME=")
	t, err := template.New("name").Parse(target)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, varmap)
	return out.String(), err
}

func loadURL(file string) (*config.Project, error) {
	resp, err := http.Get(file)
	if err != nil {
		return nil, err
	}
	p, err := config.LoadReader(resp.Body)

	// to make errcheck happy
	errc := resp.Body.Close()
	if errc != nil {
		return nil, errc
	}
	return &p, err
}

func loadFile(file string) (*config.Project, error) {
	p, err := config.Load(file)
	return &p, err
}

// Load project configuration from a given repo name or filepath/url.
func Load(repo string, file string) (project *config.Project, err error) {
	if repo == "" && file == "" {
		return nil, fmt.Errorf("Need a repo or file")
	}
	if file == "" {
		file = "https://raw.githubusercontent.com/" + repo + "/master/goreleaser.yml"
	}

	log.Printf("Reading %s", file)
	if strings.HasPrefix(file, "http") {
		project, err = loadURL(file)
	} else {
		project, err = loadFile(file)
	}
	if err != nil {
		return nil, err
	}

	// if not specified add in GitHub owner/repo info
	if project.Release.GitHub.Owner == "" {
		if repo == "" {
			return nil, fmt.Errorf("need to provide owner/name repo")
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

	return project, nil
}

func processGodownloader(repo string, filename string) {
	cfg, err := Load(repo, filename)
	if err != nil {
		log.Fatalf("Unable to parse: %s", err)
	}
	// get name template
	name, err := makeName(cfg.Archive.NameTemplate)
	cfg.Archive.NameTemplate = name
	if err != nil {
		log.Fatalf("Unable generate name: %s", err)
	}

	shell, err := makeShell(shellGodownloader, cfg)
	if err != nil {
		log.Fatalf("Unable to generate shell: %s", err)
	}
	fmt.Println(shell)
}

// processEquinoxio create a fake goreleaser config for equinox.io
// and use a similar template.
func processEquinoxio(repo string) {
	if repo == "" {
		log.Fatalf("Must have repo")
	}
	project := config.Project{}
	project.Release.GitHub.Owner = path.Dir(repo)
	project.Release.GitHub.Name = path.Base(repo)
	project.Build.Binary = path.Base(repo)
	project.Archive.Format = "tgz"

	shell, err := makeShell(shellEquinoxio, &project)
	if err != nil {
		log.Fatalf("Unable to generate shell: %s", err)
	}
	fmt.Println(shell)
}

func main() {
	source := flag.String("source", "godownloader", "download source")
	repo := flag.String("repo", "", "owner/name of repository")
	flag.Parse()
	args := flag.Args()
	file := ""
	if len(args) > 0 {
		file = args[0]
	}

	switch *source {
	case "godownloader":
		processGodownloader(*repo, file)
	case "equinoxio":
		processEquinoxio(*repo)
	default:
		log.Fatalf("Unknown source %q", *source)
	}
}
