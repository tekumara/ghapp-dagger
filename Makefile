MAKEFLAGS += --warn-undefined-variables
SHELL = /bin/bash -o pipefail
.DEFAULT_GOAL := help
.PHONY: help fmt apply state lint fix tidy test testacc vet

## display help message
help:
	@awk '/^##.*$$/,/^[~\/\.0-9a-zA-Z_-]+:/' $(MAKEFILE_LIST) | awk '!(NR%2){print $$0p}{p=$$0}' | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' | sort

## run and continuously rebuild when files change
run:
	(which reflex || go install github.com/cespare/reflex@latest)
	reflex -s -- bash -c 'source .envrc && go run .'

format:
	go fmt ./...
	cue fmt ./...
