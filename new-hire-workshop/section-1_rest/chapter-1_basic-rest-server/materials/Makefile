## help: Output this message and exit.
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'
.PHONY: help

## build: build the project
.PHONY: build
build:
	@mkdir -p dist
	go build -o dist/ ./cmd/...

## clean: clean up built binaries
.PHONY: clean
clean:
	rm -rf ./dist

## vendor: download vendored dependencies
vendor: go.mod go.sum
	go mod vendor
	@touch $@
