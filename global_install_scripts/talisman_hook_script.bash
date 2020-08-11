#!/bin/bash
shopt -s extglob
exec < /dev/tty

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

function check_and_upgrade_talisman_binary() {
	if [ -n "${TALISMAN_HOME:-}" ]; then
		LATEST_VERSION=$(curl -Is https://github.com/${ORG_REPO}/releases/latest | grep -iE "^location:" | grep -o '[^/]\+$' | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
		CURRENT_VERSION=$(${TALISMAN_BINARY} --version | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
		if [ ! -z "$LATEST_VERSION" ] && [ "$LATEST_VERSION" != "$CURRENT_VERSION" ]; then
			echo ""
			echo_warning "Your version of Talisman is outdated. Updating Talisman to v${LATEST_VERSION}"
			curl --silent https://raw.githubusercontent.com/${ORG_REPO}/master/global_install_scripts/update_talisman.bash >/tmp/update_talisman.bash && /bin/bash /tmp/update_talisman.bash talisman-binary
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

DEBUG_OPTS=""
[[ -n "${TALISMAN_DEBUG}" ]] && DEBUG_OPTS="-d"
INTERACTIVE=""
[[ -n "${TALISMAN_INTERACTIVE}" ]] && INTERACTIVE="-i"

CMD="${TALISMAN_BINARY} ${DEBUG_OPTS} --githook ${HOOKNAME} ${INTERACTIVE}"
echo_debug "ARGS are $@"
echo_debug "Executing: ${CMD}"
${CMD}
