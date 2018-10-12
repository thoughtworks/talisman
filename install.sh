#!/bin/bash
set -euo pipefail

# we call run() at the end of the script to prevent inconsistent state in case
# user runs with curl|bash and curl fails in the middle of the download
# (https://www.seancassidy.me/dont-pipe-to-your-shell.html)
run() {
  IFS=$'\n'

  VERSION="v0.3.2"
  GITHUB_URL="https://github.com/thoughtworks/talisman"
  BINARY_BASE_URL="$GITHUB_URL/releases/download/$VERSION/talisman"
  REPO_PRE_PUSH_HOOK=".git/hooks/pre-push"
  DEFAULT_GLOBAL_TEMPLATE_DIR="$HOME/.git-templates"

  EXPECTED_BINARY_SHA_LINUX_AMD64="8c0ba72fb018892b48c8e63f5e579b5bd72ec5f9d284f31c35a5382f77685834"
  EXPECTED_BINARY_SHA_LINUX_X86="332bb7a1295f45d2efaac48757f4f8c513a4cca563ebc86f964c985be7aaed51"
  EXPECTED_BINARY_SHA_DARWIN_AMD64="e66c2b21b69ab80f4474d3cc3f591f6ca68e2b76a96337e7eb807fc305e518f1"
  
  declare DOWNLOADED_BINARY
  
  E_HOOK_ALREADY_PRESENT=1
  E_CHECKSUM_MISMATCH=2
  E_USER_CANCEL=3
  E_HEADLESS=4
  E_UNSUPPORTED_ARCH=5
  E_DEPENDENCY_NOT_FOUND=6
  
  echo_error() {
    echo -ne $(tput setaf 1) >&2
    echo "$1" >&2
    echo -ne $(tput sgr0) >&2
  }

  binary_arch_suffix() {
    declare ARCHITECTURE
    if [[ "$(uname -s)" == "Linux" ]]; then
      ARCHITECTURE="linux"
    elif [[ "$(uname -s)" == "Darwin" ]]; then
      ARCHITECTURE="darwin"
    else
      echo_error "Talisman currently only supports Linux and Darwin systems."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_UNSUPPORTED_ARCH
    fi

    if [[ "$(uname -m)" = "x86_64" ]]; then
      ARCHITECTURE="${ARCHITECTURE}_amd64"
    elif [[ "$(uname -m)" =~ '^i.?86$' ]]; then
      ARCHITECTURE="${ARCHITECTURE}_386"
    else
      echo_error "Talisman currently only supports x86 and x86_64 architectures."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_UNSUPPORTED_ARCH
    fi
    
    echo $ARCHITECTURE
  }


  download_and_verify() {
    if [[ ! -x "$(which curl 2>/dev/null)" ]]; then
      echo_error "This script requires 'curl' to download the Talisman binary."
      exit $E_DEPENDENCY_NOT_FOUND
    fi
    if [[ ! -x "$(which shasum 2>/dev/null)" ]]; then
      echo_error "This script requires 'shasum' to verify the Talisman binary."
      exit $E_DEPENDENCY_NOT_FOUND
    fi
    
    echo 'Downloading and verifying binary...'
    echo
    
    TMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'talisman')
    trap 'rm -r $TMP_DIR' EXIT
    chmod 0700 $TMP_DIR

    ARCH_SUFFIX=$(binary_arch_suffix)
    
    curl --location --silent "${BINARY_BASE_URL}_${ARCH_SUFFIX}" > $TMP_DIR/talisman

    DOWNLOAD_SHA=$(shasum -b -a256 $TMP_DIR/talisman | cut -d' ' -f1)

    declare EXPECTED_BINARY_SHA
    case "$ARCH_SUFFIX" in
      linux_386)
        EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_LINUX_X86"
        ;;
      linux_amd64)
        EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_LINUX_AMD64"
        ;;
      darwin_amd64)
        EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_DARWIN_AMD64"
        ;;
    esac

    if [[ ! "$DOWNLOAD_SHA" == "$EXPECTED_BINARY_SHA" ]]; then
      echo_error "Uh oh... SHA256 checksum did not verify. Binary download must have been corrupted in some way."
      echo_error "Expected SHA: $EXPECTED_BINARY_SHA"
      echo_error "Download SHA: $DOWNLOAD_SHA"
      exit $E_CHECKSUM_MISMATCH
    fi

    DOWNLOADED_BINARY="$TMP_DIR/talisman"
  }

  install_to_repo() {
    if [[ -x "$REPO_PRE_PUSH_HOOK" ]]; then
      echo_error "Oops, it looks like you already have a pre-push hook installed at '$REPO_PRE_PUSH_HOOK'."
      echo_error "Talisman is not compatible with other hooks right now, sorry."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_HOOK_ALREADY_PRESENT
    fi

    download_and_verify

    mkdir -p $(dirname $REPO_PRE_PUSH_HOOK)
    cp $DOWNLOADED_BINARY $REPO_PRE_PUSH_HOOK
    chmod +x $REPO_PRE_PUSH_HOOK

    echo -ne $(tput setaf 2)
    echo "Talisman successfully installed to '$REPO_PRE_PUSH_HOOK'."
    echo -ne $(tput sgr0)
  }

  install_to_git_templates() {
    if [[ ! -t 1 ]]; then
      echo_error "Headless install to system templates is not supported."
      echo_error "If you would like this feature, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_HEADLESS
    fi

    TEMPLATE_DIR=$(git config --global init.templatedir) || true

    echo "Not running from inside a git repository... installing as a"
    echo "git template."
    echo
    echo "If you meant to install to a specific repo, 'cd' into that"
    echo "repo and run this script again."
    echo
    echo "Installing as a template will automatically add Talisman to"
    echo "any new repo that you 'init' or 'clone'."
    echo

    if [[ "$TEMPLATE_DIR" == "" ]]; then
      echo "No git template directory is configured. Let's add one."
      echo "(this will override any system git templates and modify your git config file)"
      echo
      read -u1 -p "Git template directory: ($DEFAULT_GLOBAL_TEMPLATE_DIR) " TEMPLATE_DIR
      echo
      TEMPLATE_DIR=${TEMPLATE_DIR:-$DEFAULT_GLOBAL_TEMPLATE_DIR}
      git config --global init.templatedir $TEMPLATE_DIR
    else
      echo "You already have a git template directory configured."
      echo
      read -u1 -p "Add Talisman to '$TEMPLATE_DIR/hooks?' (Y/n) " USE_EXISTING
      echo

      case "$USE_EXISTING" in
	Y|y|"") ;; # okay, continue
	*)
	  echo_error "Not installing Talisman."
	  echo_error "If you were trying to install into a single git repo, re-run this command from that repo."
	  echo_error "You can always download/compile manually from our Github page: $GITHUB_URL"
	  exit $E_USER_CANCEL
	  ;;
      esac
    fi

    # Support '~' in path
    TEMPLATE_DIR=${TEMPLATE_DIR/#\~/$HOME}

    if [ -f "$TEMPLATE_DIR/hooks/pre-push" ]; then
      echo_error "Oops, it looks like you already have a pre-push hook installed at '$TEMPLATE_DIR/hooks/pre-push'."
      echo_error "Talisman is not compatible with other hooks right now, sorry."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_HOOK_ALREADY_PRESENT
    fi
    
    mkdir -p "$TEMPLATE_DIR/hooks"

    download_and_verify

    cp $DOWNLOADED_BINARY "$TEMPLATE_DIR/hooks/pre-push"
    chmod +x "$TEMPLATE_DIR/hooks/pre-push"
    
    echo -ne $(tput setaf 2)
    echo "Talisman successfully installed."
    echo -ne $(tput sgr0)
  }

  if [ ! -d "./.git" ]; then
    install_to_git_templates
  else
    install_to_repo
  fi
}

run $0 $@
