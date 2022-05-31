MAKEFLAGS += --warn-undefined-variables
SHELL = /bin/bash -o pipefail
.DEFAULT_GOAL := help
.PHONY: help fmt apply state lint fix tidy test testacc vet

include .envrc

## display help message
help:
	@awk '/^##.*$$/,/^[~\/\.0-9a-zA-Z_-]+:/' $(MAKEFILE_LIST) | awk '!(NR%2){print $$0p}{p=$$0}' | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' | sort

## run and continuously rebuild when files change
run:
	(which reflex || go install github.com/cespare/reflex@latest)
	reflex -s -- bash -c 'source .envrc && go run .'

## format
format:
	go fmt ./...
	cue fmt ./...

## run dagger. Use nocache=1 to run without cache.
dagger:
	export GITHUB_TOKEN=$$(devtools/gh_token.sh $(GITHUB_REPO_URL)) && \
		dagger do ls --log-format plain $(if $(value nocache),--no-cache,)

## test
test:
	go test ./...

## debug
debug:
	@echo Connect debugger to start the app
	dlv debug --headless --listen=:12345
