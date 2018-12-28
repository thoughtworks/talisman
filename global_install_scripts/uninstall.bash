#!/bin/bash
set -euo pipefail
shopt -s extglob

DEBUG=${DEBUG:-''}

declare HOOK_SCRIPT='pre-commit' # TODO: need ability to uninstall pre-push hook as well.
if [[ $# -gt 0 && $1 =~ pre-push.* ]] ; then
   HOOK_SCRIPT='pre-push'
fi 

function run() {
    # Arguments: $1 = 'pre-commit' or 'pre-push'. whether to set talisman up as pre-commit or pre-push hook : TODO: not implemented yet
    # Environment variables:
    #    DEBUG="any-non-emply-value": verbose output for debugging the script
    #    INSTALL_ORG_REPO="..."     : the github org/repo to install from (default thoughtworks/talisman)
    #
    # Download the script needed for uninstalling the repo level hooks
    # For each git repo found in the search root (default $HOME)
    #     Run the repo level uninstall hook. This will remove the symlink from .git/hooks/pre-<commit/push> to the central $TALISMAN_SETUP_DIR
    #     The script will only remove talisman hook, not a pre-commit.com script or some other non-talisman hook
    #     Write exceptions to a file for manual action
    #     Look in the uninstall_git_repo_hook.bash script for more details on what it does
    # Remove the talisman_hook_script in .git-template/hooks/pre-<commit/push>  
    # Remove talisman binary and talisman_hook_script from $TALISMAN_SETUP_DIR ($HOME/.talisman/bin)
    
    function echo_error() {
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
    
    function echo_success {
	echo -ne $(tput setaf 2)
	echo "$1" >&2
	echo -ne $(tput sgr0)
    }
    export -f echo_success

    TALISMAN_SETUP_DIR=${HOME}/.talisman/bin
    TEMPLATE_DIR=$(git config --global init.templatedir) || true
    INSTALL_ORG_REPO=${INSTALL_ORG_REPO:-'thoughtworks/talisman'}
    SCRIPT_BASE="https://raw.githubusercontent.com/${INSTALL_ORG_REPO}/master/global_install_scripts"

    TEMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'talisman_uninstall')
    trap "rm -r ${TEMP_DIR}" EXIT
    chmod 0700 ${TEMP_DIR}

    DELETE_REPO_HOOK_SCRIPT=${TEMP_DIR}/uninstall_git_repo_hook.bash
    function get_dependent_scripts() {
	curl --silent "${SCRIPT_BASE}/uninstall_git_repo_hook.bash" > ${DELETE_REPO_HOOK_SCRIPT}
	chmod +x ${DELETE_REPO_HOOK_SCRIPT}
    }
    
    function remove_git_talisman_hooks() {
	if [[ ! -x ${DELETE_REPO_HOOK_SCRIPT} ]]; then
	    echo_error "Couldn't find executable script ${DELETE_REPO_HOOK_SCRIPT}"
	    exit 1
	fi
	
	echo "Removing talisman hooks recursively in git repos"
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
	EXCEPTIONS_FILE=${TEMP_DIR}/repos_with_multiple_hooks.paths
	touch ${EXCEPTIONS_FILE}

	TALISMAN_PATH=${TALISMAN_SETUP_DIR}/talisman_hook_script
	CMD_STRING="${SUDO_PREFIX} ${SEARCH_CMD} ${SEARCH_ROOT} ${EXTRA_SEARCH_OPTS} -name .git -type d -exec ${DELETE_REPO_HOOK_SCRIPT} ${TALISMAN_PATH} ${EXCEPTIONS_FILE} {} ${HOOK_SCRIPT} \;"
	echo_debug "EXECUTING: ${CMD_STRING}"
	eval "${CMD_STRING}" || true

	NUMBER_OF_EXCEPTION_REPOS=`cat ${EXCEPTIONS_FILE} | wc -l`

	if [ ${NUMBER_OF_EXCEPTION_REPOS} -gt 0 ]; then
	    EXCEPTIONS_FILE_HOME_PATH="${HOME}/repos_to_remove_talisman_from.paths"
	    mv ${EXCEPTIONS_FILE} ${EXCEPTIONS_FILE_HOME_PATH}
	    echo_error ""
	    echo_error "Please see ${EXCEPTIONS_FILE_HOME_PATH} for a list of repositories"
	    echo_error "that talisman couldn't be automatically removed from"
	    echo_error "This is likely because these repos are using pre-commit (https://pre-commit.com)"
	    echo_error "Remove lines related to talisman from the .pre-commit-config.yaml manually"
	fi
    }

    get_dependent_scripts
    remove_git_talisman_hooks

    echo_debug "Removing talisman hooks from .git-template"
    echo_debug "${TEMPLATE_DIR}/hooks/${HOOK_SCRIPT}"
    if [[ -n $TEMPLATE_DIR && -e ${TEMPLATE_DIR}/hooks/${HOOK_SCRIPT} && \
	      ${TALISMAN_SETUP_DIR}/talisman_hook_script -ef ${TEMPLATE_DIR}/hooks/${HOOK_SCRIPT} ]]; then
	rm -f "${TEMPLATE_DIR}/hooks/${HOOK_SCRIPT}" && \
	    echo_success "Removed ${HOOK_SCRIPT} from ${TEMPLATE_DIR}"
    fi

    echo_debug "Removing talisman from $TALISMAN_SETUP_DIR"
    rm -rf $TALISMAN_SETUP_DIR && \
	echo_success "Removed global talisman install from ${TALISMAN_SETUP_DIR}"

	if [ -n "${TALISMAN_HOME:-}" ]; then
        echo "Please remember to remove TALISMAN_HOME from your environment variables"
    fi

}

run $0 $@
