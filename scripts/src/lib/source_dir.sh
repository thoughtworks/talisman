#!/usr/bin/env bash

source_dir() {
  SOURCE="${BASH_SOURCE[0]:-$0}"
  while [ -L "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
    DIR="$(cd -P "$(dirname -- "$SOURCE")" &>/dev/null && pwd 2>/dev/null)"
    SOURCE="$(readlink -- "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="${DIR}/${SOURCE}" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
  done
  DIR="$(cd -P "$(dirname -- "$SOURCE")" &>/dev/null && pwd 2>/dev/null)"
  echo $DIR
}
