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
    CMD="git rev-list --objects --all | git cat-file --batch-check='%(objectname) %(objecttype) %(rest)' | grep \"^[^ ]* blob\" | cut -d\" \" -f1,3-"
    BINARY_PATH=${HOME}/".talisman/bin/talisman_"
    eval BLOB_DETAILS=\`${CMD}\`
    ${BINARY_PATH}${ARCHITECTURE}* -blob="${BLOB_DETAILS}"
else
    echo "Usage: talisman scan"
fi
