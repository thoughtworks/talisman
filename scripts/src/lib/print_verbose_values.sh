#!/usr/bin/env bash

function print_verbose_values_maybe() {

  if [[ -n "${args[--verbose]}" ]]; then
    echo "$(poetry --version)"
    echo "$(poetry run dbx --version)"
    echo "Databricks CLI $(poetry run databricks --version)"
  fi
}
