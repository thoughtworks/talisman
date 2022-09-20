#!/usr/bin/env bash

add_once_to_file() {
  local text=$1
  local target_file=$2

  if [[ ! -f "$target_file" ]]; then
    touch "$target_file"
  fi

  if ! grep -Fq "$text" "$target_file"
  then
    echo "$text" >>"$target_file"
  fi
}
