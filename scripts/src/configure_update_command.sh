#!/usr/bin/env bash

repo_org="${args["--repo-org"]}"
binary_name="talisman_$(__os_name)_$(__arch)"

assets=$(curl -Ls https://api.github.com/repos/"$repo_org"/releases/latest)
latest_version=$(echo "$assets" | grep tag_name | awk '{print $2}' | tr -d '"' | tr -d ',' | tr -d 'v')
current_version=$(talisman --version | awk '{print $2}')


if [[ "$latest_version" != "$current_version" ]]; then
    $(source_dir)/talisman-cli configure download
fi