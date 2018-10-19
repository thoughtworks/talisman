#!/bin/bash
set -euo pipefail
shopt -s extglob

DEBUG=${DEBUG:-''}

function run() {
    function echo_debug() {
	[[ -z "${DEBUG}" ]] && return
	echo -ne $(tput setaf 3) >&2
	echo "$1" >&2
	echo -ne $(tput sgr0) >&2
    }
    export -f echo_debug

    INSTALL_ORG_REPO=${INSTALL_ORG_REPO:-'thoughtworks/talisman'}
    SCRIPT_BASE="https://raw.githubusercontent.com/${INSTALL_ORG_REPO}/master/global_install_scripts"

    TEMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'talisman_uninstall')
    trap "rm -r ${TEMP_DIR}" EXIT
    chmod 0700 ${TEMP_DIR}

    ADD_AN_IGNORE_SCRIPT=${TEMP_DIR}/add_to_talismanignore_in_git_repo.bash

    function get_dependent_scripts() {
	echo_debug "getting ${SCRIPT_BASE}/add_to_talismanignore_in_git_repo.bash via curl"
	curl --silent "${SCRIPT_BASE}/add_to_talismanignore_in_git_repo.bash" > ${ADD_AN_IGNORE_SCRIPT}
	chmod +x ${ADD_AN_IGNORE_SCRIPT}
    }

    get_dependent_scripts
    
    echo "Adding pattern to .talismanignore recursively in git repos"
    read -p "Please enter pattern to add to .talismanignore (enter to abort): " IGNORE_PATTERN
    [[ -n $IGNORE_PATTERN ]] || exit 1
	
    read -p "Please enter root directory to search for git repos (Default: ${HOME}): " SEARCH_ROOT
    SEARCH_ROOT=${SEARCH_ROOT:-$HOME}
    SEARCH_CMD="find"
    EXTRA_SEARCH_OPTS=""
    echo -e "\tSearching ${SEARCH_ROOT} for git repositories"
    
    SUDO_PREFIX=""
    if [[ "${SEARCH_ROOT}" == "/" ]]; then
	echo -e "\tPlease enter your password when prompted to enable script to search as root user:"
	SUDO_PREFIX="sudo"
	EXTRA_SEARCH_OPTS="-xdev \( -path '/private/var' -prune \) -o"
    fi
    
    CMD_STRING="${SUDO_PREFIX} ${SEARCH_CMD} ${SEARCH_ROOT} ${EXTRA_SEARCH_OPTS} -name .git -type d -exec ${ADD_AN_IGNORE_SCRIPT} {} ${IGNORE_PATTERN} \;"
    echo_debug "EXECUTING: ${CMD_STRING}"
    eval "${CMD_STRING}"
}

run $0 $@
