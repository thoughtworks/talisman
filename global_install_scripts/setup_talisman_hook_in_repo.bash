#!/bin/bash
set -euo pipefail

function run() {
    function echo_error() {
	echo -ne $(tput setaf 1) >&2
	echo "$1" >&2
	echo -ne $(tput sgr0) >&2
    }
    
    function echo_debug() {
	[[ -z "${DEBUG}" ]] && return
	echo -ne $(tput setaf 3) >&2
	echo "$1" >&2
	echo -ne $(tput sgr0) >&2
    }
    
    function echo_success {
	echo -ne $(tput setaf 2)
	echo "$1" >&2
	echo -ne $(tput sgr0)
    }

    TALISMAN_HOOK_SCRIPT_PATH=$1
    EXCEPTIONS_FILE=$2
    DOT_GIT_DIR=$3

    REPO_HOOK_SCRIPT=${DOT_GIT_DIR}/hooks/pre-commit
    #check if a hook already exists
    if [ -e "${REPO_HOOK_SCRIPT}" ]; then
	#check if already hooked up to talisman
	if [ "${REPO_HOOK_SCRIPT}" -ef "${TALISMAN_HOOK_SCRIPT_PATH}" ]; then
	    echo_success "Talisman already setup in ${REPO_HOOK_SCRIPT}"
	else
	    if [ -e "${DOT_GIT_DIR}/../.pre-commit-config.yaml" ]; then
		echo_error "Pre-existing pre-commit.com hook detected in ${DOT_GIT_DIR}/hooks"
	    fi
	    echo ${DOT_GIT_DIR} | sed 's#/.git$##' >> ${EXCEPTIONS_FILE}
	fi
    else
	echo "Setting up pre-commit hook in ${DOT_GIT_DIR}/hooks"
	mkdir -p ${DOT_GIT_DIR}/hooks || (echo_error "Could not create hooks directory" && return)
	LN_FLAGS="-sf"
	[ -n "true" ] && LN_FLAGS="${LN_FLAGS}v"
	ln ${LN_FLAGS} ${TALISMAN_HOOK_SCRIPT_PATH} ${DOT_GIT_DIR}/hooks/pre-commit
	echo_success "DONE"
    fi
}

run $@
