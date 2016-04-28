#!/bin/bash

# we call run() at the end of the script to prevent inconsistent state in case
# curl|bash fails in the middle of the download
# (https://www.seancassidy.me/dont-pipe-to-your-shell.html)
run() {
  set -euo pipefail
  IFS=$'\n'

  GITHUB_URL="https://github.com/thoughtworks/talisman"
  BINARY_URL="$GITHUB_URL/releases/download/v0.1.0/talisman"
  EXPECTED_BINARY_SHA="fdfa31d22e5acaef3ca2f57b1036f4c2f3b9601b00524c753a5919a6c8fa3cd3"
  PRE_PUSH_HOOK=".git/hooks/pre-push"
  DOWNLOADED_BINARY=""

  echo_error() {
    echo -ne "\e[31m" >&2
    echo "$1" >&2
    echo -ne "\e[0m" >&2
  }

  download_and_verify() {
    echo
    echo 'Downloading binary...'
    
    TMP_DIR=$(mktemp -d)
    trap 'rm -r $TMP_DIR' EXIT

    curl --location --silent $BINARY_URL > $TMP_DIR/talisman

    DOWNLOAD_SHA=$(shasum -b -a256 $TMP_DIR/talisman | cut -d' ' -f1)

    if [ ! "$DOWNLOAD_SHA" = "$EXPECTED_BINARY_SHA" ]; then
      echo
      echo_error "Uh oh... SHA256 checksum did not verify. Binary download must have been corrupted in some way."
      echo_error "Expected SHA: $EXPECTED_BINARY_SHA"
      echo_error "Download SHA: $DOWNLOAD_SHA"
      exit 3
    fi

    DOWNLOADED_BINARY="$TMP_DIR/talisman"
  }

  install_to_repo() {
    if [ -x "$PRE_PUSH_HOOK" ]; then
      echo
      echo_error "Oops, it looks like you already have a pre-push hook installed at '$PRE_PUSH_HOOK'"
      echo_error "Talisman is not compatible with other hooks right now, sorry"
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit 2
    fi

    download_and_verify

    cp $DOWNLOADED_BINARY $PRE_PUSH_HOOK
    chmod +x $PRE_PUSH_HOOK

    echo
    echo -ne "\e[32m"
    echo "Talisman successfully installed to '$PRE_PUSH_HOOK'"
    echo -ne "\e[0m"
  }

  install_to_git_templates() {

    if [ ! -t 1 ]; then
      echo_error "Headless install to system templates is not supported"
      echo_error "If you would like this feature, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit 6
    fi

    echo
    echo "Talisman can be installed to your system git hook"
    echo "templates. It will then be automatically set up in any git"
    echo "repository that you 'init' or 'clone'."
    echo
    echo "This script will inform you of every file it is going to add"
    echo "to your system, and give you a chance to back out"

    echo
    read -u1 -p "Install Talisman to your system git hook templates? (Y/n) " INSTALL

    if [ "$INSTALL" != "Y" ] && [ "$INSTALL" != "y" ] && [ "$INSTALL" != "" ]; then
      echo
      echo_error "Not installing Talisman"
      echo_error "If you were trying to install into a single git repo, re-run this command from that repo"
      echo_error "You can always download/compile manually from our Github page: $GITHUB_URL"
      exit 4
    fi

    echo
    echo "Searching for your git template directories..."
    GIT_HOOK_TEMPLATE_DIRS=$(find / -path '*/git-core/templates/hooks' 2>/dev/null) || true

    echo
    echo "Installing Talisman as a 'pre-push' binary in the following locations:"
    echo
    echo $GIT_HOOK_TEMPLATE_DIRS
    echo
    read -u1 -p "Continue? (Y/n) " CONTINUE
    
    if [ "$CONTINUE" != "Y" ] && [ "$CONTINUE" != "y" ] && [ "$CONTINUE" != "" ]; then
      echo
      echo_error "Not installing Talisman"
      echo_error "You can always download/compile manually from our Github page: $GITHUB_URL"
      exit 5
    fi

    download_and_verify

    echo
    echo "Installing Talisman to the above locations (may ask for sudo access)"
    for DIR in "$GIT_HOOK_TEMPLATE_DIRS"; do
      sudo cp $DOWNLOADED_BINARY "$DIR/pre-push"
      sudo chmod +x "$DIR/pre-push"
    done
    
    echo
    echo -ne "\e[32m"
    echo "Talisman successfully installed"
    echo -ne "\e[0m"
  }

  if [ ! -d "./.git" ]; then
    install_to_git_templates
  else
    install_to_repo
  fi
}

run $0 $@
