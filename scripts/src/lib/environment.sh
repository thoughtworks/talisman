#!/usr/bin/env bash

function __osx() {
  # __osx - return if we are running in 'Darwin' (macOS)
  [[ $(uname) == "Darwin" ]]
}

function __ubuntu() {
  [[ $(uname) == "Linux" && $(cat /etc/issue | grep -i ubuntu) ]]
}

function noop() {
  echo 'noop' >/dev/null
}

function __os_name() {
  local my_os
  my_os=$(uname -s)
  echo "${my_os,,}"
}

function __arch() {
  local my_arch
  my_arch=$(uname -m)

  local suffix
  case $my_arch in
  "x86_64")
    suffix="amd64"
    ;;
  "i686" | "i386")
    suffix="386"
    ;;
  "arm64")
    suffix="arm64"
    ;;
  *)
    echo_error "Unsupported architecture: $my_arch"
    exit 1
    ;;
  esac

  echo "$suffix"
}
