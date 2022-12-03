#!/bin/bash
set -euo pipefail

GIT_REPO_DOT_GIT=$1
IGNORE_PATTERN=$2

function echo_debug() {
  [[ -z "${DEBUG}" ]] && return
  echo -ne $(tput setaf 3) >&2
  echo "$1" >&2
  echo -ne $(tput sgr0) >&2
}

function echo_success {
  echo -ne $(tput setaf 2)
  echo "$1" >&2
  echo -ne $(tput sgr0)
}

echo_debug "Adding $IGNORE_PATTERN to ${GIT_REPO_DOT_GIT}/../.talismanignore"
echo $IGNORE_PATTERN >>${GIT_REPO_DOT_GIT}/../.talismanignore && echo_success "Added $IGNORE_PATTERN to ${GIT_REPO_DOT_GIT}/../.talismanignore"
