SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=

install: build ## build and install
	go install .

setup: ## Install all the build and lint dependencies
	./scripts/setup.sh

test: ## Run all the tests
	gotestcover $(TEST_OPTIONS) -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	./scripts/lint.sh

precommit:  ## Run precommit hook
	./scripts/lint.sh

ci: build lint test  ## travis-ci entrypoint
	./samples/godownloader-goreleaser.sh
	git diff .
	./bin/goreleaser --snapshot

build: install_hooks ## Build a beta version of goreleaser
	go build
	./scripts/build_samples.sh

.DEFAULT_GOAL := build

generate: ## regenerate shell code from client9/shlib
	./makeshellfn.sh > shellfn.go

.PHONY: ci help generate samples clean

clean: ## clean up everything
	go clean ./...
	rm -f godownloader
	rm -rf ./bin ./dist
	git gc --aggressive

install_hooks:  ## install precommit hooks for git
	cp -f scripts/lint.sh .git/hooks/pre-commit

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

