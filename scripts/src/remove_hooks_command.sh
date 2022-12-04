#!/usr/bin/env bash

hook_script_path="${args["--hook-script-path"]}"
exceptions_file="${args["--exceptions-file-path"]}"
git_dir="${args["--git-dir"]}"
hook_name="${args["--hook-name"]}"

cyan_ln "Removing Talisman $hook_name hook in $git_dir"


local_hook_script_path=${git_dir}/hooks/${hook_name}

magenta_ln "Removing $local_hook_script_path"

if [[ ! -e "${local_hook_script_path}" ]]; then
   gray_ln "File ${local_hook_script_path} does not exist, nothing to do"
  exit 0
fi

rm -f "$local_hook_script_path"

# remove script if it is symlinked to talisman
if [ "${REPO_HOOK_SCRIPT}" -ef "${hook_script_path}" ]; then
  rm ${REPO_HOOK_SCRIPT} && cyan_ln "Removed ${REPO_HOOK_SCRIPT}"
  exit 0
fi

if [ -e "${git_dir}/../.pre-commit-config.yaml" ]; then
  # check if the .pre-commit-config contains "talisman", if so ask them to remove it manually
  yellow_ln "Pre-existing pre-commit.com hook detected in ${git_dir}/hooks"
fi

echo ${git_dir} | sed 's#/.git$##' >>$exceptions_file
