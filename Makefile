SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
OS=$(shell uname -s)

export PATH := ./bin:$(PATH)

setup: ## Install all the build and lint dependencies
	mkdir -p bin
	go get -u golang.org/x/tools/cmd/cover
	curl -sfL https://install.goreleaser.com/github.com/gohugoio/hugo.sh | bash
	curl -sfL https://install.goreleaser.com/github.com/alecthomas/gometalinter.sh | bash
ifeq ($(OS), Darwin)
	brew install dep
else
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure
	gometalinter --install
.PHONY: setup

install: build ## build and install
	go install .

test: ## Run all the tests
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	gometalinter --vendor ./...

precommit: lint  ## Run precommit hook

ci: build lint test  ## travis-ci entrypoint
	./samples/godownloader-goreleaser.sh
	git diff .
	./bin/goreleaser --snapshot

build: hooks ## Build a beta version of goreleaser
	go build
	./scripts/build-site.sh

.DEFAULT_GOAL := build

generate: ## regenerate shell code from client9/shlib
	./makeshellfn.sh > shellfn.go

.PHONY: ci help generate samples clean

clean: ## clean up everything
	go clean ./...
	rm -f godownloader
	rm -rf ./bin ./dist ./vendor
	git gc --aggressive

hooks:
	echo "make lint" > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

