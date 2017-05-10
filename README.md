# godownloader
Download Go binaries as fast and easily as possible.

This is the inverse of [goreleaser](https://github.com/goreleaser/goreleaser).  The goreleaser YAML file is read and creates a custom shell script that can download the right package and the right version for the existing machine.  (sometimes this is known as "curl bash")

This is useful in CI/CD systems such as [travis-ci.org](https://travis-ci.org).

* Much faster then 'go get' (sometimes up to 100x)
* Make sure your local environment (macOS) and CI enviroment (Linux) are using the exact same versions of your go binaries.

## Example

Let's say you are using [hugo](https://gohugo.io), the static website generator, with [travis-ci](https://travis-ci.org).

Your old `.travis.yml` file might have 

```
install:
  - go get github.com/spf13/hugo
```

This can take up to 30 seconds!  To fix this up first do:

```
# create a godownloader script
godownloader -repo spf/hugo > ./hugo-installer.sh`
```

and add `hugo-installer.sh` to your GitHub repo.  Edit your `.travis.yml` as such

```yaml
install:
  - ./hugo-install.sh 0.20.6
```

There is even experimental download latest function:

```yaml
install:
  - ./hugo-install.sh latest
```

Typical download time is 0.3 seconds, or 100x improvement.

Your new `hugo` binary is in `./bin`, so change your Makefie or scripts to use `./bin/hugo`. 

## Status

* Supports GitHub Releases only
* Binary is installed using `tar.gz` (`zip` support is coming)
* Windows anything.  I just don't know enough about it.
* OS-specific install such as homebrew, deb, rpm aren't supported.  Everything is installed locally via .tar.gz.

## TODO

* Zip support and format over-rides.
* Should `VERSION` be set via ENV not arg `$1`
* Checksum support
* Add LICENSE
* Adjustment of default `TMPDIR` and `BINDIR` (install directory)
* Vendor dependencies
* Figure out how to write tests
* Add a ton of comments
* Setup travis.ci
* Use goreleaser to release godownloader
* Use godownloader to download godownloader

## Yes, it's true.

It's a go program that reads a YAML file that uses a template to make a posix shell script.

