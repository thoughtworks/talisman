#!/usr/bin/env bash

repo_org="${args["--repo-org"]}"
binary_name="talisman_$(__os_name)_$(__arch)"

hook_script_path="${args["--hook-script-path"]}"
exceptions_file="${args["--exceptions-file-path"]}"
git_dir="${args["--git-dir"]}"
hook_name="${args["--hook-name"]}"

white_ln "Installing talisman in $git_dir"

git_hook_script_path=${git_dir}/hooks/${hook_name}

#check if a hook already exists
if [ -e "${git_hook_script_path}" ]; then
  #check if already hooked up to talisman
  if [ "${git_hook_script_path}" -ef "${hook_script_path}" ]; then
    green_ln "Talisman already setup in ${git_hook_script_path}"
  else
    if [ -e "${DOT_GIT_DIR}/../.pre-commit-config.yaml" ]; then
      red_ln "Pre-existing pre-commit.com hook detected in ${git_dir}/hooks"
    fi
    yellow_ln ${git_dir} | sed 's#/.git$##' >>${exceptions_file}
  fi
else
  yellow_ln "Setting up ${hook_script_path} hook in ${git_dir}/hooks"
  mkdir -p ${git_dir}/hooks || (echo_error "Could not create hooks directory" && return)
  LN_FLAGS="-sf"

  [ -n "true" ] && LN_FLAGS="${LN_FLAGS}v"

  ln ${LN_FLAGS} "${hook_script_path}" "${git_dir}"/hooks/"${hook_name}"

  green_ln "DONE"
fi

# TODO : Add windows support
