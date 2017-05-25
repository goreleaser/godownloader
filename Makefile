SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=

setup: ## Install all the build and lint dependencies
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/...
	go get -u github.com/pierrre/gotestcover
	go get -u golang.org/x/tools/cmd/cover
	dep ensure
	gometalinter --install

test: ## Run all the tests
	gotestcover $(TEST_OPTIONS) -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s
	shellcheck samples/godownloader-goreleaser.sh

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...

lint_shell:  ## shellcheck the shell scripts
	shellcheck -s sh samples/godownloader-goreleaser.sh
	shellcheck -s bash samples/godownloader-goreleaser.sh
	shellcheck -s dash samples/godownloader-goreleaser.sh
	shellcheck -s ksh samples/godownloader-goreleaser.sh

ci: build samples lint test lint_shell ## Run all the tests and code checks as travis-ci does
	./samples/godownloader-goreleaser.sh
	./bin/goreleaser --snapshot

build: ## Build a beta version of goreleaser
	go build

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build

generate: ## regenerate shell code from client9/posixshell
	./makeshellfn.sh > shellfn.go

samples: ## make sample donwloaders
	./godownloader -repo spf13/hugo > samples/godownloader-hugo.sh
	./godownloader -repo goreleaser/goreleaser > samples/godownloader-goreleaser.sh
	./godownloader -repo client9/misspell > samples/godownloader-misspell.sh$
	./godownloader -source equinoxio -repo tdewolff/minify > samples/godownloader-minify.sh
	chmod a+x samples/*.sh

.PHONY: ci help generate samples clean

clean: ## clean up everything
	go clean ./...
	rm -f godownloader
	rm -rf ./bin ./dist
	git gc --aggressive

