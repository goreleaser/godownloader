# go-downloader
The inverse of go-releaser.

reads a gorelease.yaml file and generates a posix shell
script to download it that works on any OS (with exception of windows).


Useful in CI/CD systems such as travis-ci.org

* Much faster then 'go get'
* Make sure your local environment (macOS) and the CI are using the exact same versions.

## Yes, it's true.

It's a go program that reads a YAML file that uses a template to make a bash file.


