#!/bin/bash
# Hello there! If you update the talisman version, please remember to:
# - Test that this script works with no args, and the `pre-push` / `pre-commit` arg.
# - also update `install.sh` in the gh_pages branch of this repo, so that
#   <https://thoughtworks.github.io/talisman/install.sh> gets updated too.
# Thanks!

set -euo pipefail

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
  IFS=$'\n'

  VERSION="v1.8.0"
  GITHUB_URL="https://github.com/thoughtworks/talisman"
  BINARY_BASE_URL="$GITHUB_URL/releases/download/$VERSION/talisman"
  REPO_HOOK_BIN_DIR=".git/hooks/bin"

  DEFAULT_GLOBAL_TEMPLATE_DIR="$HOME/.git-templates"

  EXPECTED_BINARY_SHA_LINUX_AMD64="22b1aaee860b27306bdf345a0670f138830bcf7fbe16c75be186fe119e9d54b4"
  EXPECTED_BINARY_SHA_LINUX_X86="d0558d626a4ee1e90d2c2a5f3c69372a30b8f2c8e390a59cedc15585b0731bc4"
  EXPECTED_BINARY_SHA_DARWIN_AMD64="f30e1ec6fb3e1fc33928622f17d6a96933ca63d5ab322f9ba869044a3075ffda"
  EXPECTED_BINARY_SHA_WINDOWS_X86="98ee5ed4bb394096a643531b7b8d3e6e919cc56e4673add744b46036260527c3"
  EXPECTED_BINARY_SHA_WINDOWS_AMD64="697cebb5988ee002b630b814c6c6f5d49d921c9c3aad4545c4a77d749e5ae833"

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

  echo_success() {
    echo -ne $(tput setaf 2)
    echo "$1" >&2
    echo -ne $(tput sgr0)
  }

  binary_arch_suffix() {
    declare ARCHITECTURE
    UNAME_S="$(uname -s)"
    UNAME_O="$(uname -o)"
    if [[ "$UNAME_S" == "Linux" ]]; then
      ARCHITECTURE="linux"
    elif [[ "$UNAME_S" == "Darwin" ]]; then
      ARCHITECTURE="darwin"
    elif [[ "$UNAME_O" == "Msys" ]]; then
      ARCHITECTURE="windows"
    else
      echo_error "Talisman currently only supports Linux, Windows (Git Bash) and Darwin systems."
      echo_error "OS Detected: ${UNAME_S}, ${UNAME_O}."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_UNSUPPORTED_ARCH
    fi

    if [[ "$(uname -m)" = "x86_64" ]]; then
      ARCHITECTURE="${ARCHITECTURE}_amd64"
    elif [[ "$(uname -m)" =~ ^i.?86$ ]]; then
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
    chmod 0700 "$TMP_DIR"

    ARCH_SUFFIX=$(binary_arch_suffix)

    declare DOWNLOAD_SUFFIX
    if [[ "$ARCH_SUFFIX" =~ ^windows ]]; then
      DOWNLOAD_SUFFIX="${ARCH_SUFFIX}.exe"
    else
      DOWNLOAD_SUFFIX="$ARCH_SUFFIX"
    fi;

    curl --location --silent "${BINARY_BASE_URL}_${DOWNLOAD_SUFFIX}" >"${TMP_DIR}/talisman"

    DOWNLOAD_SHA=$(shasum -b -a256 "${TMP_DIR}/talisman" | cut -d' ' -f1)

    declare EXPECTED_BINARY_SHA
    case "$ARCH_SUFFIX" in
    linux_386)
      EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_LINUX_X86"
      ;;
    linux_amd64)
      EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_LINUX_AMD64"
      ;;
    windows_386)
      EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_WINDOWS_X86"
      ;;
    windows_amd64)
      EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_WINDOWS_AMD64"
      ;;
    darwin_amd64)
      EXPECTED_BINARY_SHA="$EXPECTED_BINARY_SHA_DARWIN_AMD64"
      ;;
    esac

    if [[ ! "$DOWNLOAD_SHA" == "$EXPECTED_BINARY_SHA" ]]; then
      echo_error "Uh oh... SHA256 checksum did not verify. Binary download must have been corrupted in some way."
      echo_error "Arch Suffix: ${ARCH_SUFFIX}"
      echo_error "Expected SHA: ${EXPECTED_BINARY_SHA}"
      echo_error "Download SHA: ${DOWNLOAD_SHA}"
      echo_error "Downloaded binary: ${TMP_DIR}/talisman"
      exit $E_CHECKSUM_MISMATCH
    fi

    DOWNLOADED_BINARY="$TMP_DIR/talisman"
  }

  install_to_repo() {
    if [[ -x "$REPO_HOOK_TARGET" ]]; then
      echo_error "Oops, it looks like you already have a ${HOOK_NAME} hook installed at '${REPO_HOOK_TARGET}'."
      echo_error "Talisman is not compatible with other hooks right now, sorry."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
      exit $E_HOOK_ALREADY_PRESENT
    fi

    download_and_verify

    mkdir -p "$REPO_HOOK_BIN_DIR"
    TALISMAN_BIN_TARGET="${REPO_HOOK_BIN_DIR}/talisman"
    cp "$DOWNLOADED_BINARY" "$TALISMAN_BIN_TARGET"
    chmod +x "$TALISMAN_BIN_TARGET"

    cat >"$REPO_HOOK_TARGET" <<EOF
#!/bin/bash
TALISMAN_OPTS="--githook ${HOOK_NAME}"
[[ -n "\${TALISMAN_DEBUG}" ]] && TALISMAN_OPTS="-d \${TALISMAN_OPTS}"
TALISMAN_BIN="${PWD}/${TALISMAN_BIN_TARGET}"
[[ -n "\${TALISMAN_DEBUG}" ]] && echo "ARGS are \$@"
[[ -n "\${TALISMAN_DEBUG}" ]] && echo "Executing: \${TALISMAN_BIN} \${TALISMAN_OPTS}"
"\$TALISMAN_BIN" \$TALISMAN_OPTS
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
      echo_error "Talisman is not compatible with other hooks right now, sorry."
      echo_error "If this is a problem for you, please open an issue: https://github.com/thoughtworks/talisman/issues/new"
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
