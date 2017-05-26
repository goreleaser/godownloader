# godownloader
Download Go binaries as fast and easily as possible

This is the inverse of [goreleaser](https://github.com/goreleaser/goreleaser).  The goreleaser YAML file is read and creates a custom shell script that can download the right package and the right version for the existing machine.

If you use goreleaser already, this will create scripts suitable for "curl bash" style downloads.

* Run godownloader on your `goreleaser.yaml` file
* Add the `godownloader.sh` file to your repo
* Tell your users to use https://raw.githubusercontent.com/YOU/YOURAPP/master/godownloader.sh to install

This is also useful in CI/CD systems such as [travis-ci.org](https://travis-ci.org).

* Much faster then 'go get' (sometimes up to 100x)
* Make sure your local environment (macOS) and CI environment (Linux) are using the exact same versions of your go binaries.

## CI/CD Example

Let's say you are using [hugo](https://gohugo.io), the static website generator, with [travis-ci](https://travis-ci.org).

Your old `.travis.yml` file might have 

```yaml
install:
  - go get github.com/spf13/hugo
```

This can take up to 30 seconds! 

Hugo doesn't have (yet) a `godownloader.sh` file.  So we will make our own:


```
# create a godownloader script
godownloader -repo spf/hugo > ./godownloader-hugo.sh`
```

and add `godownloader-hugo.sh` to your GitHub repo.  Edit your `.travis.yml` as such

```yaml
install:
  - ./godownloader-hugo.sh 0.20.6
```

Without a version number, GitHub is queried to get the latest version number.  This is subject to the usual [GitHub rate limits](https://developer.github.com/v3/#rate-limiting).  If working on a public machine (like travis-ci), be sure to set `GITHUB_TOKEN`.

```yaml
install:
  - ./godownloader-hugo.sh
```

Typical download time is 0.3 seconds, or 100x improvement. 

Your new `hugo` binary is in `./bin`, so change your Makefie or scripts to use `./bin/hugo`. 

The default installation directory can be changed with the `-b` flag.

## Notes on Functionality

* Only GitHub Releases are supported right now.
* Checksums are checked.
* Binares are installed using `tar.gz` or `zip`. 
* No support for Windows anything.  I just don't know enough about it.
* No OS-specific installs such as homebrew, deb, rpm.  Everything is installed locally via a `tar.gz` or `zip`.  Typically OS installs are done differently anyways (e.g. brew, apt-get, yum, etc).

## Experimental support

Some people do not use Goreleaser!  

There is experimental support for the following alterative distributions.

### "naked" releases on GitHub

A naked release is just the raw binary put on GitHub releases.  Limited support can be done by

```bash
./goreleaser -source raw -repo [owner/repo] -exe [name] -nametpl [tpl]
```

Where `exe` is the final binary name, and `tpl` is the same type of name template that Goreleaser uses.

An example repo is at [mvdan/sh](https://github.com/mvdan/sh/releases). Note how repo `sh` is different than binary `shfmt`.

### Equinox.io

[Equinox.io](https://equinox.io) is a really interesting platform.  Take a look.

There is no API, so godownloader scripts screen scrapes to figure out the latest release.  Likewise, checksums are not verified.

```bash
./goreleaser -source equinoxio -repo [owner/repo]
```

While Equinox.io supports the concepts of channels, it is hardwired to `stable` for now.

An example is [tdewolff/minify](https://github.com/tdewolff/minify) on [dl.equinox.io](https://dl.equinox.io/tdewolff/minify/stable).

## Yes, it's true.

It's a go program that reads a YAML file that uses a template to make a posix shell script.

