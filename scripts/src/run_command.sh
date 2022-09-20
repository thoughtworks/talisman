#!/usr/bin/env bash


#!/bin/bash
shopt -s extglob

# Download and install updates if available

# Don't run talisman checks in a git repo, if we find a .talisman_skip or .talisman_skip.pre-<commit/push> file in the repo
# Run interactive mode if we are running pre-commit

#CMD="${TALISMAN_BINARY} ${DEBUG_OPTS} --githook ${HOOKNAME} ${INTERACTIVE}"
#${CMD}
