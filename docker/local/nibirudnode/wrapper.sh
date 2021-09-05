#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/nibirud/${BINARY:-nibirud}
ID=${ID:-0}
LOG=${LOG:-nibirud.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'nibirud' E.g.: -e BINARY=nibirud_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export NIBIRUDHOME="/nibirud/node${ID}/nibirud"

if [ -d "$(dirname "${NIBIRUDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${NIBIRUDHOME}" "$@" | tee "${NIBIRUDHOME}/${LOG}"
else
  "${BINARY}" --home "${NIBIRUDHOME}" "$@"
fi
