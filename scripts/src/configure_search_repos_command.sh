#!/usr/bin/env bash
path="${args["--path"]}"

echo -e "\tSearching ${path} for git repositories"

EXTRA_SEARCH_OPTS=""
SUDO_PREFIX=""
if [[ "${SEARCH_ROOT}" == "/" ]]; then
  echo -e "\tPlease enter your password when prompted to enable script to search as root user:"
  SUDO_PREFIX="sudo"
  EXTRA_SEARCH_OPTS="-xdev \( -path '/private/var' -prune \) -o"
fi

CMD_STRING="${SUDO_PREFIX} find ${path} ${EXTRA_SEARCH_OPTS} -name .git -type d"

eval "${CMD_STRING}" || true
