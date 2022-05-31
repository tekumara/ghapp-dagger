#!/bin/bash

# Usage:
#       gh_token github_host
#
# eg: source devtools/gh_token.sh && export GITHUB_TOKEN=$(gh_token github.com)

# get the oauth token from the github cli config
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

    if ! hash yq; then
        echo -e "https://github.com/mikefarah/yq is not installed.\nPlease install, eg: brew install yq" >&2
        return 42
    fi

    if ! token=$(yq e ".\"$github_host\".oauth_token" -P ~/.config/gh/hosts.yml) || [[ $token == null ]]; then
        echo -e "No token found for $github_host.\nPlease login into github, eg: gh auth login -h $github_host" >&2
        return 42
    fi

    echo "$token"
}

gh_token "$1"
