#!/bin/bash
#
# Usage: vault-write <path> ["secret strings>" | @<secret file>]

SCRIPT_NAME=$(basename "$0")
VAULT_ADDR=${VAULT_ADDR:-http//localhost:8200}
if [ "$#" != "2" ]; then
    echo usage "$SCRIPT_NAME: <path> [\"<secret strings>\" | @<secret file>]"
    exit 1
fi

path=$1

if [[ "$2" =~ ^@ ]];
then
    # if the data is in a file
    src=$(echo $2 | cut -c 2-)
    if file -b --mime-encoding $src | grep -s binary > /dev/null
    then
        # if data is binary, base64 encode it and set format=base64
        cat $src | base64 | vault kv put $path value=- format="base64"
    else
        # otherwise set format=text
        cat $src | vault kv put $path value=- format="text"
    fi
else
    vault kv put $path value="$2" format="text"
fi
