#!/bin/bash
shopt -s extglob

# set TALISMAN_DEBUG="some-non-empty-value" in the env to get verbose output when the hook or talisman is running
function echo_debug() {
    MSG="$@"
    [[ -n "${TALISMAN_DEBUG}" ]] && echo "${MSG}"
}

declare HOOKNAME="pre-commit"
NAME=$(basename $0)

# given the various symlinks, this script may be invoked as
#     'pre-commit', 'pre-push', 'talisman_hook_script pre-commit' or 'talisman_hook_script pre-push'
case "$NAME" in
    pre-commit*|pre-push*) HOOKNAME="${NAME}" ;;
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

# Here HOOKNAME should be either 'pre-commit' (default) or 'pre-push'
echo_debug "Firing ${HOOKNAME} hook"

# Don't run talisman checks in a git repo, if we find a .talisman_skip or .talisman_skip.pre-<commit/push> file in the repo
if [[ -f .talisman_skip || -f .talisman_skip.${HOOKNAME} ]] ; then
	echo_debug "Found skip file. Not performing checks"
	exit 0
fi

DEBUG_OPTS=""
[[ -n "${TALISMAN_DEBUG}" ]] && DEBUG_OPTS="-d"
CMD="${TALISMAN_BINARY} ${DEBUG_OPTS} -githook ${HOOKNAME}"
echo_debug "ARGS are $@"
echo_debug "Executing: ${CMD}"
${CMD}
