#!/bin/bash
# Hello there! If you update the talisman version, please remember to:
# - Test that this script works with no args, and the `pre-push` / `pre-commit` arg.
# - also update `install.sh` in the gh_pages branch of this repo, so that
#   <https://thoughtworks.github.io/talisman/install.sh> gets updated too.
# Thanks!

set -euo pipefail

DEBUG=${DEBUG:-''}
HOOK_NAME="${1:-pre-push}"
case "$HOOK_NAME" in
pre-commit | pre-push) REPO_HOOK_TARGET=".git/hooks/${HOOK_NAME}" ;;
*)
  echo "Unknown Hook name '${HOOK_NAME}'. Please check parameters"
  exit 1
  ;;
esac

# we call run() at the end of the script to prevent inconsistent state in case
# user runs with curl|bash and curl fails in the middle of the download
# (https://www.seancassidy.me/dont-pipe-to-your-shell.html)
run() {
  declare TALISMAN_BINARY_NAME

  IFS=$'\n'

  GITHUB_URL="https://github.com/thoughtworks/talisman"
  VERSION=$(curl --silent "https://api.github.com/repos/thoughtworks/talisman/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  BINARY_BASE_URL="$GITHUB_URL/releases/download/$VERSION"
  REPO_HOOK_BIN_DIR=".git/hooks/bin"

  DEFAULT_GLOBAL_TEMPLATE_DIR="$HOME/.git-templates"

  declare DOWNLOADED_BINARY
  TEMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'talisman_setup')
	trap "rm -r ${TEMP_DIR}" EXIT
	chmod 0700 ${TEMP_DIR}

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
  export -f echo_error

	function echo_debug() {
		[[ -z "${DEBUG}" ]] && return
		echo -ne $(tput setaf 3) >&2
		echo "$1" >&2
		echo -ne $(tput sgr0) >&2
	}
	export -f echo_debug

  echo_success() {
    echo -ne $(tput setaf 2)
    echo "$1" >&2
    echo -ne $(tput sgr0)
  }
  export -f echo_success

  operating_system() {
    OS=$(uname -s)
    case $OS in
      "Linux")
        echo "linux"
        ;;
      "Darwin")
        echo "darwin"
        ;;
      MINGW32_NT-10.0-WOW*)
        echo "windows"
        ;;
      MINGW64_NT-10.0*)
        echo "windows"
        ;;
      *)
        echo_error "Talisman currently only supports Windows, Linux and MacOS(darwin) systems."
        echo_error "If this is a problem for you, please open an issue: https://github.com/${INSTALL_ORG_REPO}/issues/new"
        exit $E_UNSUPPORTED_ARCH
        ;;
      esac
}

   binary_arch_suffix() {
    declare OS
    OS=$(operating_system)
		ARCH=$(uname -m)
		case $ARCH in
		"x86_64")
			OS="${OS}_amd64"
			;;
		"i686" | "i386")
			OS="${OS}_386"
			;;
		*)
			echo_error "Talisman currently only supports x86 and x86_64 architectures."
			echo_error "If this is a problem for you, please open an issue: https://github.com/${INSTALL_ORG_REPO}/issues/new"
			exit $E_UNSUPPORTED_ARCH
			;;
		esac

		TALISMAN_BINARY_NAME="talisman_${OS}"
		if [[ $OS == *"windows"* ]]; then
			TALISMAN_BINARY_NAME="${TALISMAN_BINARY_NAME}.exe"
		fi
  }

	function download() {
		OBJECT=$1
		DOWNLOAD_URL=${BINARY_BASE_URL}/${OBJECT}
		echo "Downloading ${OBJECT} from ${DOWNLOAD_URL}"
		curl --location --silent ${DOWNLOAD_URL} >"$TEMP_DIR/${OBJECT}"
	}

  function verify_checksum() {
		FILE_NAME=$1
		CHECKSUM_FILE_NAME='checksums'
		echo_debug "Verifying checksum for ${FILE_NAME}"
		download ${CHECKSUM_FILE_NAME}

		pushd ${TEMP_DIR} >/dev/null 2>&1
		grep ${TALISMAN_BINARY_NAME} ${CHECKSUM_FILE_NAME} >${CHECKSUM_FILE_NAME}.single
		shasum -a 256 -c ${CHECKSUM_FILE_NAME}.single
		popd >/dev/null 2>&1
		echo_debug "Checksum verification successful!"
		echo
	}

	function download_and_verify() {
		binary_arch_suffix
		download "${TALISMAN_BINARY_NAME}"
		DOWNLOADED_BINARY="${TEMP_DIR}/${TALISMAN_BINARY_NAME}"
		verify_checksum "${TALISMAN_BINARY_NAME}"
	}

  install_to_repo() {
    if [[ -x "$REPO_HOOK_TARGET" ]]; then
      echo_error "Oops, it looks like you already have a ${HOOK_NAME} hook installed at '${REPO_HOOK_TARGET}'."
      echo_error "If this is expected, you should consider setting-up a tool to allow git hook chaining,"
      echo_error "like pre-commit (brew install pre-commit) or Husky or any other tool of your choice."
      echo_error "WARNING! Talisman hook not installed."
      exit $E_HOOK_ALREADY_PRESENT
    fi

    download_and_verify

    mkdir -p "$REPO_HOOK_BIN_DIR"
    TALISMAN_BIN_TARGET="${REPO_HOOK_BIN_DIR}/talisman"
    cp "$DOWNLOADED_BINARY" "$TALISMAN_BIN_TARGET"
    chmod +x "$TALISMAN_BIN_TARGET"

    cat >"$REPO_HOOK_TARGET" <<EOF
#!/bin/bash
[[ -n "\${TALISMAN_DEBUG}" ]] && DEBUG_OPTS="-d"
CMD="${PWD}/${TALISMAN_BIN_TARGET} \${DEBUG_OPTS} --githook ${HOOK_NAME}"
[[ -n "\${TALISMAN_DEBUG}" ]] && echo "ARGS are \$@"
[[ -n "\${TALISMAN_DEBUG}" ]] && echo "Executing: \${CMD}"
\${CMD}
EOF
    chmod +x "$REPO_HOOK_TARGET"

    echo_success "Talisman successfully installed to '$REPO_HOOK_TARGET'."
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
      Y | y | "") ;; # okay, continue
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

    if [ -f "$TEMPLATE_DIR/hooks/${HOOK_NAME}" ]; then
      echo_error "Oops, it looks like you already have a ${HOOK_NAME} hook installed at '$TEMPLATE_DIR/hooks/${HOOK_NAME}'."
				echo_error "If this is expected, you should consider setting-up a tool to allow git hook chaining,"
				echo_error "like pre-commit (brew install pre-commit) or Husky or any other tool of your choice."
				echo_error "WARNING! Talisman hook not installed."
      exit $E_HOOK_ALREADY_PRESENT
    fi

    mkdir -p "$TEMPLATE_DIR/hooks"

    download_and_verify

    cp "$DOWNLOADED_BINARY" "$TEMPLATE_DIR/hooks/${HOOK_NAME}"
    chmod +x "$TEMPLATE_DIR/hooks/${HOOK_NAME}"

    echo_success "Talisman successfully installed."
  }

  if [ ! -d "./.git" ]; then
    install_to_git_templates
  else
    install_to_repo
  fi
}

run $0 $@
