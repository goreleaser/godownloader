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
	"time"

	"github.com/goreleaser/goreleaser/config"
	"github.com/goreleaser/goreleaser/context"
	"github.com/goreleaser/goreleaser/pipeline/defaults"
)

// given a template, and a config, generate shell script
func makeShell(tplsrc string, cfg *config.Project) (string, error) {

	// if we want to add a timestamp in the templates this
	//  function will generate it
	funcMap := template.FuncMap{
		"timestamp": func() string {
			return time.Now().UTC().Format(time.RFC3339)
		},
	}

	var out bytes.Buffer
	t, err := template.New("shell").Funcs(funcMap).Parse(tplsrc)
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
// {{ .Binary }} --->  [prefix]${BINARY}, etc.
//
func makeName(prefix, target string) (string, error) {
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
		"Os":          "${OS}",
		"Arch":        "${ARCH}",
		"Version":     "${VERSION}",
		"Tag":         "${TAG}",
		"Binary":      "${BINARY}",
		"ProjectName": "${PROJECT_NAME}",
	}

	var out bytes.Buffer
	out.WriteString(prefix)
	t, err := template.New("name").Parse(target)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, varmap)
	return out.String(), err
}

func loadURLs(path string) (*config.Project, error) {
	for _, file := range []string{"goreleaser.yml", ".goreleaser.yml", "goreleaser.yaml", ".goreleaser.yaml"} {
		var url = fmt.Sprintf("%s/%s", path, file)
		log.Printf("reading %s", url)
		project, err := loadURL(url)
		if err != nil {
			return nil, err
		}
		if project != nil {
			return project, nil
		}
	}
	return nil, fmt.Errorf("could not fetch a goreleaser configuration file")
}

func loadURL(file string) (*config.Project, error) {
	resp, err := http.Get(file)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("reading %s returned %d %s\n", file, resp.StatusCode, http.StatusText(resp.StatusCode))
		return nil, nil
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
		return nil, fmt.Errorf("repo or file not specified")
	}
	if file == "" {
		log.Printf("read repo %q on github", repo)
		project, err = loadURLs(
			fmt.Sprintf("https://raw.githubusercontent.com/%s/master", repo),
		)
	} else {
		log.Printf("reading file %q", file)
		project, err = loadFile(file)
	}
	if err != nil {
		return nil, err
	}

	// if not specified add in GitHub owner/repo info
	if project.Release.GitHub.Owner == "" {
		if repo == "" {
			return nil, fmt.Errorf("owner/name repo not specified")
		}
		project.Release.GitHub.Owner = path.Dir(repo)
		project.Release.GitHub.Name = path.Base(repo)
	}

	var ctx = context.New(*project)
	err = defaults.Pipe{}.Run(ctx)
	project = &ctx.Config

	// set default binary name
	if len(project.Builds) == 0 {
		project.Builds = []config.Build{
			{Binary: path.Base(repo)},
		}
	}
	if project.Builds[0].Binary == "" {
		project.Builds[0].Binary = path.Base(repo)
	}

	return project, err
}

func main() {
	var (
		source  = flag.String("source", "godownloader", "download source")
		exe     = flag.String("exe", "", "name of binary, used only in raw")
		nametpl = flag.String("nametpl", "", "name template, used only in raw")
		repo    = flag.String("repo", "", "owner/name of repository")
	)

	flag.Parse()
	args := flag.Args()
	file := ""
	if len(args) > 0 {
		file = args[0]
	}
	var (
		out string
		err error
	)
	switch *source {
	case "godownloader":
		// https://github.com/goreleaser/godownloader
		out, err = processGodownloader(*repo, file)
	case "equinoxio":
		// https://equinox.io
		out, err = processEquinoxio(*repo)
	case "raw":
		// raw mode is when people upload direct binaries
		// to GitHub releases that are not  not tar'ed or zip'ed.
		// For example:
		//   https://github.com/mvdan/sh/releases
		out, err = processRaw(*repo, *exe, *nametpl)
	default:
		log.Fatalf("unknown source %q", *source)
	}

	if err != nil {
		log.Fatalf("failed: %s", err)
	}
	fmt.Print(out)
}
