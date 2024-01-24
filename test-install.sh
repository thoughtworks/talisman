_linux_uname() {
  if [ "${FAKE_PARAMS[0]}" = "-m" ]; then
    echo "x86_64"
  else
    echo "Linux"
  fi
}
export -f _linux_uname

_windows_uname() {
  if [ "${FAKE_PARAMS[0]}" = "-m" ]; then
    echo "i686"
  else
    echo "MINGW64_NT-10.0-19045"
  fi
}
export -f _windows_uname

_mac_uname() {
  if [ "${FAKE_PARAMS[0]}" = "-m" ]; then
    echo "aarch64"
  else
    echo "Darwin"
  fi
}
export -f _mac_uname

setup() {
  temp=$(mktemp -d)
  fake uname _linux_uname
  fake curl 'echo "download_url: talisman_linux_amd64checksums"'
  fake shasum true
  fake tput true
}

teardown() {
  rm -rf "$temp"
}

test_installs_without_sudo() {
  fake sudo 'echo "expected no sudo" && exit 1'
  INSTALL_LOCATION=$temp ./install.sh
  assert "test -x $temp/talisman_linux_amd64" "Should install file with executable mode"
  assert_matches "$temp/talisman_linux_amd64" "$(readlink "$temp/talisman")" "Should create a link"
}

test_installs_with_sudo_if_available() {
  fake touch 'echo "Permission denied" && exit 1'
  fake which 'echo "sudo installed" && exit 0'
  # shellcheck disable=SC2016
  fake sudo 'bash -c "${FAKE_PARAMS[*]}"'
  INSTALL_LOCATION=$temp ./install.sh
  assert "test -x $temp/talisman_linux_amd64" "Should install file with executable mode"
  assert_matches "$temp/talisman_linux_amd64" "$(readlink "$temp/talisman")" "Should create a link"
}

test_errors_if_unable_to_install() {
  fake touch 'echo "Permission denied" && exit 1'
  fake which 'echo "sudo not installed" && exit 1'
  assert_status_code 126 "INSTALL_LOCATION=$temp ./install.sh"
}

test_errors_if_no_install_location() {
  assert_status_code 1 "INSTALL_LOCATION=/does/not/exist ./install.sh"
}

test_mac_arm_binary_name() {
  fake uname _mac_uname
  fake curl 'echo "download_url: talisman_darwin_arm64checksums"'
  INSTALL_LOCATION=$temp ./install.sh
  assert "test -x $temp/talisman_darwin_arm64" "Should install file with executable mode"
  assert_matches "$temp/talisman_darwin_arm64" "$(readlink "$temp/talisman")" "Should create a link"
}

test_windows_binary_name() {
  fake uname _windows_uname
  fake curl 'echo "download_url: talisman_windows_386.exechecksums"'
  INSTALL_LOCATION=$temp ./install.sh
  assert "test -x $temp/talisman_windows_386.exe" "Should install file with executable mode"
  assert_matches "$temp/talisman_windows_386.exe" "$(readlink "$temp/talisman")" "Should create a link"
}
