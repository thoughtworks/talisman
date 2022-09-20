#!/usr/bin/env bash

set -euo pipefail

function _redact_string() {
  local target_string
  target_string=$1

  local num_visisble_chars
  num_visisble_chars=$2

  local count
  count=$((${#target_string}-$num_visisble_chars))

  local stars
  stars=$(printf '*%.0s' echo $(eval echo {0..$count}))

  local is_left
  is_left=$3

  local sed_str
  if [[ $is_left == "true" ]]; then
    sed_str="s/^.\{${count}\}/${stars}/g"
  else
    sed_str="s/.\{${count}\}$/${stars}/g"
  fi

  target_string_masked=$(echo $target_string | sed -e $sed_str)
  echo $target_string_masked
}

function redact_string_right() {
  _redact_string $1 $2 false
}

function redact_string_left() {
  _redact_string $1 $2 true
}

set +euo pipefail
