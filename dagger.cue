package main

import (
	"dagger.io/dagger"
	"universe.dagger.io/alpine"
	"universe.dagger.io/bash"
	"universe.dagger.io/docker"
	"universe.dagger.io/git"
)

dagger.#Plan & {
	client: env: {
		GITHUB_TOKEN:    dagger.#Secret
		GITHUB_REPO_URL: string
		GITHUB_REF:      string
	}

	actions: {
		_pull: git.#Pull & {
			remote: client.env.GITHUB_REPO_URL
			ref:    client.env.GITHUB_REF
			auth: authToken: client.env.GITHUB_TOKEN
		}

		_alpine: alpine.#Build & {
			packages: bash: _
		}

		_copy: docker.#Copy & {
			input:    _alpine.output
			contents: _pull.output
		}

		ls: bash.#Run & {
			input: _copy.output
			script: contents: "ls -al"
		}
	}
}
