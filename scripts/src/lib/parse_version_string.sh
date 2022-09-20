#!/usr/bin/env bash

function get_version_major() {

  IFS='.' read -ra version_arr <<< "$1"
  echo "${version_arr[0]}"
}

function get_version_minor() {

  IFS='.' read -ra version_arr <<< "$1"
  echo "${version_arr[1]}"
}

function get_version_patch() {

  IFS='.' read -ra version_arr <<< "$1"
  echo "${version_arr[2]}"
}
