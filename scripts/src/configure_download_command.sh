#!/usr/bin/env bash

repo_org="${args["--repo-org"]}"
binary_name="talisman_$(__os_name)_$(__arch)"

assets=$(curl -Ls https://api.github.com/repos/"$repo_org"/releases/latest)
download_url=$(echo "$assets" | grep download_url | awk '{print $2}' | tr -d '"' | grep "$binary_name")
checksum_url=$(echo "$assets" | grep download_url | awk '{print $2}' | tr -d '"' | grep "checksum")

temp_dir=$(mktemp -d)

cyan_ln "Downloading talisman binary."
mkdir -p $HOME/.talisman/bin
curl --location --silent "${download_url}" >"$temp_dir"/"$binary_name"
curl --location --silent "${checksum_url}" | grep "$binary_name" >"$temp_dir"/checksums

pushd "$temp_dir" 2>&1 >/dev/null || exit
sha256sum -c checksums

$(source_dir)/talisman-cli configure download-cli

if [ $? -eq 0 ]; then
  mv "$temp_dir"/"$binary_name" "$HOME"/.talisman/bin/talisman
  chmod +x "$HOME"/.talisman/bin/talisman
  cp "$HOME"/.talisman/bin/talisman "$HOME"/.talisman/bin/"$binary_name"
  rm -rf "$temp_dir"
  green_ln "Talisman binary downloaded successfully."
  popd 2>&1 >/dev/null || exit
else
  red_ln "Talisman binary download failed."
  popd 2>&1 >/dev/null || exit
  exit 1
fi
