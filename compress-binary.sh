#!/usr/bin/env bash

set -e

pushd dist || exit

binary=$1

if [[ $binary == *"darwin"* -o  $binary == *"arm"* ]]; then
  exit 0
fi

echo "Compressing binaries"
upx --lzma "$binary"
echo "...Done"

echo "Verifying compressed binaries"
upx -t "$binary"
echo "...Done"

popd || exit
