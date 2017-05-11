# godownloader
Download Go binaries as fast and easily as possible

This is the inverse of [goreleaser](https://github.com/goreleaser/goreleaser).  The goreleaser YAML file is read and creates a custom shell script that can download the right package and the right version for the existing machine.

If you use goreleaser already, this will create script suitable for "curl bash" style downloads.

* Run godownloader on your `goreleaser.yaml` file
* Add the `godownloader.sh` file to your repo
* Tell your users to use https://raw.githubusercontent.com/YOU/YOURAPP/master/godownloader.sh to install

This is also useful in CI/CD systems such as [travis-ci.org](https://travis-ci.org).

* Much faster then 'go get' (sometimes up to 100x)
* Make sure your local environment (macOS) and CI enviroment (Linux) are using the exact same versions of your go binaries.

## CI/CD Example

Let's say you are using [hugo](https://gohugo.io), the static website generator, with [travis-ci](https://travis-ci.org).

Your old `.travis.yml` file might have 

```yaml
install:
  - go get github.com/spf13/hugo
```

This can take up to 30 seconds! 

Hugo doesn't have (yet!) a godownloader.sh file.  So we will make our own:


```
# create a godownloader script
godownloader -repo spf/hugo > ./godownloader-hugo.sh`
```

and add `godownloader-hugo.sh` to your GitHub repo.  Edit your `.travis.yml` as such

```yaml
install:
  - ./godownloader-hugo.sh 0.20.6
```

There is even experimental download latest function:

```yaml
install:
  - ./godownloader-hugo.sh latest
```

Typical download time is 0.3 seconds, or 100x improvement.

Your new `hugo` binary is in `./bin`, so change your Makefie or scripts to use `./bin/hugo`. 

## Status

* Only GitHub Releases are supported right now.
* Binares are installed using `tar.gz` or `zip`. 
* No support for Windows anything.  I just don't know enough about it.
* No OS-specific installs such as homebrew, deb, rpm.  Everything is installed locally via a `tar.gz` or `zip`.  Typically OS installs are done differently anyways (e.g. yum, apt-get, etc).

## TODO

* #5 Checksum support
* #8 Setup travis.ci
* #10 Adjustment of default `BINDIR` (install directory)
* #11 Use goreleaser to release godownloader
* #12 Use godownloader to download godownloader
* Vendor dependencies
* Figure out how to write tests
* Add a ton of comments
* Should `VERSION` be set via ENV not arg `$1` ?

## Yes, it's true.

It's a go program that reads a YAML file that uses a template to make a posix shell script.

