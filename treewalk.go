package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/client9/codegen/shell"
	"gopkg.in/yaml.v2"
)

type treeConfig struct {
	// these can be set by config
	Source  string `yaml:"source,omitempty"`  // type of downloader to make
	Exe     string `yaml:"exe,omitempty"`     // stuff for "raw"
	Nametpl string `yaml:"nametpl,omitempty"` // stuff for "raw"

	// these can not be set by config file
	// and are set by the url/path
	org   string // github.com for now
	owner string // ^ github username
	name  string // repo name
}

// Load config file
func LoadTreeConfig(file string) (config treeConfig, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	log.WithField("file", file).Debug("loading config file")
	return LoadTreeConfigReader(f)
}

// LoadReader config via io.Reader
func LoadTreeConfigReader(fd io.Reader) (config treeConfig, err error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return config, err
	}
	err = yaml.UnmarshalStrict(data, &config)
	log.WithField("config", config).Debug("loaded config file")
	return config, err
}

// for many files, this might be slow since golang reads and sorts
// everything.  If it's a problem, investigate:
//
func treewalk(root string, treeout string, forceWrite bool) error {
	rooterr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// weird case where filewalk failed
		if err != nil {
			return err
		}

		// ignore directories
		if info.IsDir() {
			return nil
		}
		suffix := filepath.Ext(path)
		// ignore non-yaml stuff
		if suffix != ".yaml" && suffix != ".yml" {
			return nil
		}

		// Now: root/github.com/owner/repo.yaml
		rel, err := filepath.Rel(root, path)
		if err != nil {
			panic("should never happen.. path is always in root")
		}

		// Now: github.com/owner/repo.yaml
		rel = rel[0 : len(rel)-len(suffix)]

		// Now: github.com/owner/repo
		// better way of doing this?
		parts := strings.Split(rel, string(os.PathSeparator))
		if len(parts) != 3 {
			return fmt.Errorf("invalid path: %s", path)
		}

		org, owner, repo := parts[0], parts[1], parts[2]
		// Now: [ github.com client misspell ]

		// only github.com for now
		if org != "github.com" {
			return fmt.Errorf("only github.com supported, got %s", org)
		}

		// nice and clean
		//  org == github.com
		//  owner == you
		//  repo  == your project

		c, err := LoadTreeConfig(path)
		if err != nil {
			return err
		}

		// hacking for now and just hardwiring
		if c.Source == "" {
			c.Source = "godownloader"
		}

		// overwrite what exists for security
		c.org = org
		c.owner = owner
		c.name = repo

		shellcode, err := processSource(c.Source, owner+"/"+repo, "", c.Exe, c.Nametpl)
		if err != nil {
			return err
		}

		// now write back
		outdir := filepath.Join(treeout, org, repo)
		err = os.MkdirAll(outdir, 0755)
		if err != nil {
			return err
		}
		shellpath := filepath.Join(outdir, owner+".sh")

		// only write out if forced to, OR if output is effectively different
		// than what the file has.
		if forceWrite || shell.ShouldWriteFile(shellpath, shellcode) {
			if err = ioutil.WriteFile(shellpath, shellcode, 0755); err != nil {
				return err
			}
		}

		// we did it!
		return nil
	})
	return rooterr
}
