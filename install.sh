#!/bin/bash
set -euo pipefail

# we call run() at the end of the script to prevent inconsistent state in case
# user runs with curl|bash and curl fails in the middle of the download
# (https://www.seancassidy.me/dont-pipe-to-your-shell.html)

if [ $# -eq 0 ]
  then
    echo "Installing the pre-push hook which is default, to install the pre-commit hook please provide the "pre-commit" parameter to install-talisman shell script"
    REPO_HOOK=".git/hooks/pre-push"
    HOOK="pre-push"
fi

if [ $# -eq 1 ]
  then
    if [ $1 == 'pre-push' ]
       then 
          REPO_HOOK=".git/hooks/pre-push"
          HOOK="pre-push"
    elif [ $1 == 'pre-commit' ]
         then
           REPO_HOOK=".git/hooks/pre-commit"
           HOOK="pre-commit"
    else
      echo "Talisman only supports either pre-push or pre-commit hooks"
      exit 0
    fi
else
  echo "The number of arguements provided to the script are incorrect"
  exit 0
fi

run() {
  IFS=$'\n'
  VERSION="v0.2.0"
  GITHUB_URL="https://github.com/karanmilan/talisman"
  BINARY_BASE_URL="$GITHUB_URL/releases/download/$VERSION/talisman"
  DEFAULT_GLOBAL_TEMPLATE_DIR="$HOME/.git-templates"
  
  EXPECTED_BINARY_SHA_LINUX_AMD64="92b9a75e04451ca71cb7b75dae395cced4b68e3d62a0df68d7062fda3b32563d"
  EXPECTED_BINARY_SHA_LINUX_X86="181f2f88153be912bfc835563d55bea2b5b83c259312f489a71ad3a0358c097e"
  EXPECTED_BINARY_SHA_DARWIN_AMD64="f5084c41f2350423f39fd2e719397e959f306e5e18ab0daff4be69547106f936"
  
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


  download() {
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
    
    curl --location --silent "${BINARY_BASE_URL}_${ARCH_SUFFIX}_${HOOK}" > $TMP_DIR/talisman

    DOWNLOADED_BINARY="$TMP_DIR/talisman"
  }

  install_to_repo() {
          
    if [[ -x "$REPO_HOOK" ]]; then
      echo_error "Oops, it looks like you already have a similar hook installed at '$REPO_HOOK'."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_HOOK_ALREADY_PRESENT
    fi

    download

    mkdir -p $(dirname $REPO_HOOK)
    cp $DOWNLOADED_BINARY $REPO_HOOK
    chmod +x $REPO_HOOK

    echo -ne $(tput setaf 2)
    echo "Talisman successfully installed to '$REPO_HOOK'."
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

    download

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
