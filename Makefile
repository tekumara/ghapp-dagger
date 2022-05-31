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

## format
format:
	go fmt ./...
	cue fmt ./...

## run dagger. Use nocache=1 to run without cache.
GITHUB_HOST ?= github.com
GITHUB_REPO_URL ?= https://github.com/tekumara/ghapp-dagger.git
GITHUB_REF ?= dd142c005ca6e2decf64b7e295e0f858d6195820
dagger:
	source devtools/gh_token.sh && export GITHUB_TOKEN=$$(gh_token $(GITHUB_HOST)) && \
		dagger do build --log-format plain $(if $(value nocache),--no-cache,)
