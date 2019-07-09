#!/bin/bash
set -euo pipefail
shopt -s extglob

DEBUG=${DEBUG:-''}
FORCE_DOWNLOAD=${FORCE_DOWNLOAD:-''}
declare UPDATE_TYPE=""
if [[ $# -gt 0 && $1 =~ talisman-binary.* ]]; then
	UPDATE_TYPE='talisman-binary'
fi

function run() {

	# Download appropriate (appropriate = based on OS and ARCH) talisman binary from github
	# Copy the talisman binary to $TALISMAN_SETUP_DIR ($HOME/.talisman/bin)

	declare TALISMAN_BINARY_NAME

	IFS=$'\n'
	VERSION=${VERSION:-'latest'}
	INSTALL_ORG_REPO=${INSTALL_ORG_REPO:-'thoughtworks/talisman'}

	TALISMAN_SETUP_DIR=${HOME}/.talisman/bin # location of central install: talisman binary and hook script

	TEMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'talisman_setup')
	#trap "rm -r ${TEMP_DIR}" EXIT
	chmod 0700 ${TEMP_DIR}

	function echo_error() {
		echo -ne $(tput setaf 1) >&2
		echo "$1" >&2
		echo -ne $(tput sgr0) >&2
	}
	export -f echo_error

	function echo_debug() {
		[[ -z "${DEBUG}" ]] && return
		echo -ne $(tput setaf 3) >&2
		echo "$1" >&2
		echo -ne $(tput sgr0) >&2
	}
	export -f echo_debug

	function echo_success() {
		echo -ne $(tput setaf 2)
		echo "$1" >&2
		echo -ne $(tput sgr0)
	}
	export -f echo_success

	function collect_version_artifact_download_urls() {
		curl -Ls -w %{url_effective} "https://github.com/${INSTALL_ORG_REPO}/releases/latest" | grep -Eo '/'${INSTALL_ORG_REPO}'/releases/download/.+/[^/"]+' | sed 's/^/https:\/\/github.com/' >${TEMP_DIR}/download_urls
		echo_debug "All release artifact download urls can be found at ${TEMP_DIR}/download_urls:"
		[[ -z "${DEBUG}" ]] && return
		cat ${TEMP_DIR}/download_urls
	}

	function set_talisman_binary_name() {
		# based on OS (linux/darwin) and ARCH(32/64 bit)
		echo_debug "Running set_talisman_binary_name"
		declare ARCHITECTURE
		OS=$(uname -s)
		case $OS in
		"Linux")
			ARCHITECTURE="linux"
			;;
		"Darwin")
			ARCHITECTURE="darwin"
			;;
		"MINGW32_NT-10.0-WOW")
			ARCHITECTURE="windows"
			;;
		"MINGW64_NT-10.0")
			ARCHITECTURE="windows"
			;;
		*)
			echo_error "Talisman currently only supports Windows, Linux and MacOS(darwin) systems."
			echo_error "If this is a problem for you, please open an issue: https://github.com/${INSTALL_ORG_REPO}/issues/new"
			exit $E_UNSUPPORTED_ARCH
			;;
		esac

		ARCH=$(uname -m)
		case $ARCH in
		"x86_64")
			ARCHITECTURE="${ARCHITECTURE}_amd64"
			;;
		"i686" | "i386")
			ARCHITECTURE="${ARCHITECTURE}_386"
			;;
		*)
			echo_error "Talisman currently only supports x86 and x86_64 architectures."
			echo_error "If this is a problem for you, please open an issue: https://github.com/${INSTALL_ORG_REPO}/issues/new"
			exit $E_UNSUPPORTED_ARCH
			;;
		esac

		TALISMAN_BINARY_NAME="talisman_${ARCHITECTURE}"
		if [[ "$OS" == "MINGW32_NT-10.0-WOW" || "$OS" == "MINGW64_NT-10.0" ]]; then
			TALISMAN_BINARY_NAME="${TALISMAN_BINARY_NAME}.exe"
		fi
	}

	function download() {
		echo_debug "Running download()"
		OBJECT=$1
		DOWNLOAD_URL=$(grep 'http.*'${OBJECT}'$' ${TEMP_DIR}/download_urls)
		echo_debug "Downloading ${OBJECT} from ${DOWNLOAD_URL}"
		curl --location --silent ${DOWNLOAD_URL} >${TEMP_DIR}/${OBJECT}
	}

	function verify_checksum() {
		FILE_NAME=$1
		CHECKSUM_FILE_NAME='checksums'
		echo_debug "Verifying checksum for ${FILE_NAME}"
		download ${CHECKSUM_FILE_NAME}

		pushd ${TEMP_DIR} >/dev/null 2>&1
		grep ${TALISMAN_BINARY_NAME} ${CHECKSUM_FILE_NAME} >${CHECKSUM_FILE_NAME}.single
		shasum -a 256 -c ${CHECKSUM_FILE_NAME}.single
		popd >/dev/null 2>&1
		echo_debug "Checksum verification successful!"
		echo
	}

	function download_talisman_binary() {
		#download talisman binary
		echo_debug "Running download_talisman_binary"
		download ${TALISMAN_BINARY_NAME}
		verify_checksum ${TALISMAN_BINARY_NAME}
	}

	function download_talisman_hook_script() {
		echo_debug "Running download_talisman_hook_script"
		curl --silent https://raw.githubusercontent.com/${INSTALL_ORG_REPO}/master/global_install_scripts/talisman_hook_script.bash >${TEMP_DIR}/talisman_hook_script
	}

	function setup_talisman() {
		# copy talisman binary from TEMP folder to the central location
		rm -f ${TALISMAN_SETUP_DIR}/${TALISMAN_BINARY_NAME}
		cp ${TEMP_DIR}/${TALISMAN_BINARY_NAME} ${TALISMAN_SETUP_DIR}
		chmod +x ${TALISMAN_SETUP_DIR}/${TALISMAN_BINARY_NAME}
		echo_success "Talisman binary updated successfully!"
	}

	function setup_talisman_hook_script() {
		BINARY_PATH=${TALISMAN_SETUP_DIR}/${TALISMAN_BINARY_NAME}
		rm -f ${TALISMAN_SETUP_DIR}/talisman_hook_script
		sed "s@\${TALISMAN_BINARY}@$BINARY_PATH@g" ${TEMP_DIR}/talisman_hook_script >${TALISMAN_SETUP_DIR}/talisman_hook_script
		chmod +x ${TALISMAN_SETUP_DIR}/talisman_hook_script
		echo_success "Talisman hook script updated successfully!"
	}

	function set_talisman_env_variables_properly() {
		FILE_PATH="$1"
		if [ -f $FILE_PATH ] && grep -q "TALISMAN_HOME" $FILE_PATH; then
			if ! grep -q ">>> talisman >>>" $FILE_PATH; then
				sed -i'-talisman.bak' '/TALISMAN_HOME/d' $FILE_PATH
				echo -e "\n" >>${ENV_FILE}
				echo "# >>> talisman >>>" >>$FILE_PATH
				echo "# Below environment variables should not be modified unless you know what you are doing" >>$FILE_PATH
				echo "export TALISMAN_HOME=${TALISMAN_SETUP_DIR}" >>$FILE_PATH
				echo "alias talisman=\$TALISMAN_HOME/${TALISMAN_BINARY_NAME}" >>$FILE_PATH
				echo "# <<< talisman <<<" >>$FILE_PATH
			fi
		fi
	}

	set_talisman_binary_name
	echo "Downloading latest talisman binary..."
	collect_version_artifact_download_urls
	download_talisman_binary
	setup_talisman
	if [ -z "$UPDATE_TYPE" ]; then
		echo "Downloading latest talisman hook script..."
		download_talisman_hook_script
		setup_talisman_hook_script
	fi
	# Correcting talisman env variables if they are not in proper format
	if [ -n "${TALISMAN_HOME:-}" ]; then
		set_talisman_env_variables_properly ~/.bashrc
		set_talisman_env_variables_properly ~/.bash_profile
		set_talisman_env_variables_properly ~/.profile
	fi

}

run $0 $@
