#!/usr/bin/env bash

repo_org="${args["--repo-org"]}"
github_url="https://raw.githubusercontent.com"

mkdir -p ${HOME}/.talisman/bin/

#curl --silent "$github_url/${repo_org}/javier-installation-refactor/talisman-cli" >${HOME}/.talisman/bin/talisman-cli
cat "/Users/javiermatias-cabrera/workspace/external-repos/talisman/talisman-cli" >${HOME}/.talisman/bin/talisman-cli
chmod +x ${HOME}/.talisman/bin/talisman-cli

green_ln "Installed talisman-cli version $(${HOME}/.talisman/bin/talisman-cli -v) at ${HOME}/.talisman/bin/talisman-cli"
