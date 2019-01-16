#!/bin/bash

COMMAND_PASSED=$1
if [ ${COMMAND_PASSED} = "scan" ]; then
    	declare ARCHITECTURE
	OS=$(uname -s)
	case $OS in
	    "Linux")
		ARCHITECTURE="linux" ;;
	    "Darwin")
		ARCHITECTURE="darwin" ;;
		"MINGW32_NT-10.0-WOW")
		ARCHITECTURE="windows" ;;
		"MINGW64_NT-10.0")
		ARCHITECTURE="windows" ;;
	esac
	COMMITS_COMMAND="git log --reflog --pretty=oneline | cut -d \" \" -f1"
	eval COMMITS=\`${COMMITS_COMMAND}\`
	BLOBS=""
	for COMMIT in ${COMMITS}
	do
		BLOBS_CMD="git ls-tree -r ${COMMIT} | cut -d \" \" -f3,4"
		eval CURRENT_BLOBS=\`${BLOBS_CMD}\`
		BLOBS="${BLOBS}commit ${COMMIT} ${CURRENT_BLOBS}"
	done
    BINARY_PATH=${HOME}/".talisman/bin/talisman_"
    ${BINARY_PATH}${ARCHITECTURE}* -blob="${BLOBS}"
else
    echo "Usage: talisman scan"
fi
