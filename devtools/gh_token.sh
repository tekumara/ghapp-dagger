#!/bin/bash

# Get the oauth token from the github cli config
#
# Usage:
#       gh_token github_host
#
# eg:
#       export GITHUB_TOKEN=$(gh_token.sh github.com)

gh_token() {
    local github_host=${1:-}
    local token

    if [[ -z "${1:-}" ]]; then
        echo -e "Missing github host, eg: $0 github.com" >&2
        return 42
    elif [[ "$1" =~ https?://.* ]]; then
        # extract hostname from url
        github_host=$(echo "$1" | awk -F'[/:]' '{print $4}')
    else
        github_host=$1
    fi

    if ! hash gh; then
        echo -e "https://github.com/cli/cli is not installed.\nPlease install, eg: brew install gh" >&2
        return 42
    fi

    gh auth token -h "$github_host"

    echo "$token"
}

gh_token "$1"
