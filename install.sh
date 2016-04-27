#!/bin/bash

# we call run() at the end of the script to prevent inconsistent state in case
# curl|bash fails in the middle of the download
# (https://www.seancassidy.me/dont-pipe-to-your-shell.html)
run() {
  set -euo pipefail

  BINARY_URL="https://github.com/thoughtworks/talisman/releases/download/v0.1.0/talisman"
  EXPECTED_BINARY_SHA="fdfa31d22e5acaef3ca2f57b1036f4c2f3b9601b00524c753a5919a6c8fa3cd3"
  PRE_PUSH_HOOK=".git/hooks/pre-push"

  echo_error() {
    echo -ne "\e[31m" >&2
    echo "$1" >&2
    echo -ne "\e[0m" >&2
  }

  if [ ! -d "./.git" ]; then
    echo_error "Oops, please run me from the root of the git repo that you want to add Talisman to"
    exit 1
  fi

  if [ -x "$PRE_PUSH_HOOK" ]; then
    echo_error "Oops, it looks like you already have a pre-push hook installed at '$PRE_PUSH_HOOK'"
    echo_error "Talisman is not compatible with other hooks right now, sorry"
    echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
    exit 2
  fi

  TMP_DIR=$(mktemp -d)
  trap 'rm -r $TMP_DIR' EXIT

  curl --location --silent $BINARY_URL > $TMP_DIR/talisman

  DOWNLOAD_SHA=$(shasum -b -a256 $TMP_DIR/talisman | cut -d' ' -f1)

  if [ ! "$DOWNLOAD_SHA" = "$EXPECTED_BINARY_SHA" ]; then
    echo_error "Uh oh... SHA256 checksum did not verify. Binary download must have been corrupted in some way."
    echo_error "Expected SHA: $EXPECTED_BINARY_SHA"
    echo_error "Download SHA: $DOWNLOAD_SHA"
    exit 3
  fi

  cp $TMP_DIR/talisman $PRE_PUSH_HOOK
  chmod +x $PRE_PUSH_HOOK

  echo -ne "\e[32m"
  echo "Talisman successfully installed to '$PRE_PUSH_HOOK'"
  echo -ne "\e[0m"
}

run $0 $@
