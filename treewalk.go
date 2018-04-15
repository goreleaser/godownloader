package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/client9/codegen/shell"
)

type sourceConfig struct {
	source  string // type of downloader to make
	org     string // github.com for now
	owner   string // ^ github username
	repo    string // repo name
	exe     string // stuff for "raw"
	nametpl string // stuff for "raw"
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
		if suffix != ".yaml" {
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

		log.Printf("PATH: %s %s %s", org, owner, repo)

		c := sourceConfig{}

		// REAL PATH AS YAML

		// hacking for now and just hardwiring
		c.source = "godownloader"

		// overwrite what exists for security
		c.org = org
		c.owner = owner
		c.repo = repo

		shellcode, err := processSource(c.source, owner+"/"+repo, "", c.exe, c.nametpl)
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
