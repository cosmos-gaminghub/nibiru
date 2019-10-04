#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/nbrd/${BINARY:-nbrd}
ID=${ID:-0}
LOG=${LOG:-nbrd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'nbrd' E.g.: -e BINARY=nbrd_my_test_version"
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
export NBRDHOME="/nbrd/node${ID}/nbrd"

if [ -d "`dirname ${NBRDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$NBRDHOME" "$@" | tee "${NBRDHOME}/${LOG}"
else
  "$BINARY" --home "$NBRDHOME" "$@"
fi

chmod 777 -R /nbrd

