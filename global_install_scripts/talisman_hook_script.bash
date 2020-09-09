#!/bin/bash
shopt -s extglob

# set TALISMAN_DEBUG="some-non-empty-value" in the env to get verbose output when the hook or talisman is running
function echo_debug() {
	MSG="$@"
	[[ -n "${TALISMAN_DEBUG}" ]] && echo "${MSG}"
}

function echo_warning() {
	echo -ne $(tput setaf 3) >&2
	echo "$1" >&2
	echo -ne $(tput sgr0) >&2
}

function echo_error() {
	echo -ne $(tput setaf 1) >&2
	echo "$1" >&2
	echo -ne $(tput sgr0) >&2
}

function echo_success() {
	echo -ne $(tput setaf 2)
	echo "$1" >&2
	echo -ne $(tput sgr0)
}

function toLower(){
	echo "$1" | awk '{print tolower($0)}'
}

declare HOOKNAME="pre-commit"
NAME=$(basename $0)
ORG_REPO=${ORG_REPO:-'thoughtworks/talisman'}

# given the various symlinks, this script may be invoked as
#     'pre-commit', 'pre-push', 'talisman_hook_script pre-commit' or 'talisman_hook_script pre-push'
case "$NAME" in
pre-commit* | pre-push*) HOOKNAME="${NAME}" ;;
talisman_hook_script)
	if [[ $# -gt 0 && $1 =~ pre-push.* ]]; then
		HOOKNAME="pre-push"
	fi
	;;
*)
	echo "Unexpected invocation. Please check invocation name and parameters"
	exit 1
	;;
esac

TALISMAN_UPGRADE_CONNECT_TIMEOUT=${TALISMAN_UPGRADE_CONNECT_TIMEOUT:-10}
function check_and_upgrade_talisman_binary() {
	if [[ -n "${TALISMAN_HOME:-}" && "$TALISMAN_SKIP_UPGRADE" != "true" ]]; then
		LATEST_VERSION=$(curl --connect-timeout $TALISMAN_UPGRADE_CONNECT_TIMEOUT -Is https://github.com/${ORG_REPO}/releases/latest | grep -iE "^location:" | grep -o '[^/]\+$' | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
		CURRENT_VERSION=$(${TALISMAN_BINARY} --version | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
		if [ -z "$LATEST_VERSION" ]; then
			echo_warning "Failed to retrieve latest version, skipping update."
		elif [ "$LATEST_VERSION" != "$CURRENT_VERSION" ]; then
			echo ""
			echo_warning "Your version of Talisman is outdated. Updating Talisman to v${LATEST_VERSION}"
			curl --silent https://raw.githubusercontent.com/${ORG_REPO}/master/global_install_scripts/update_talisman.bash >/tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash
		else
			echo_debug "Talisman version up-to-date, skipping update"
		fi
	fi
}

check_and_upgrade_talisman_binary
# Here HOOKNAME should be either 'pre-commit' (default) or 'pre-push'
echo_debug "Firing ${HOOKNAME} hook"

# Don't run talisman checks in a git repo, if we find a .talisman_skip or .talisman_skip.pre-<commit/push> file in the repo
if [[ -f .talisman_skip || -f .talisman_skip.${HOOKNAME} ]]; then
	echo_debug "Found skip file. Not performing checks"
	exit 0
fi

TALISMAN_DEBUG="$(toLower "${TALISMAN_DEBUG}")"
DEBUG_OPTS=""
[[ "${TALISMAN_DEBUG}" == "true" ]] && DEBUG_OPTS="-d"

TALISMAN_INTERACTIVE="$(toLower "${TALISMAN_INTERACTIVE}")"
INTERACTIVE=""
if [ "${TALISMAN_INTERACTIVE}" == "true" ]; then
    INTERACTIVE="-i"
	  [[ "${HOOKNAME}" == "pre-commit" ]] && exec < /dev/tty || echo_warning "talisman pre-push hook cannot be invoked in interactive mode currently"
fi

CMD="${TALISMAN_BINARY} ${DEBUG_OPTS} --githook ${HOOKNAME} ${INTERACTIVE}"
echo_debug "ARGS are $@"
echo_debug "Executing: ${CMD}"
${CMD}
