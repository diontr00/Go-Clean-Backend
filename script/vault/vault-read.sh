#!/bin/bash
# Usage: vault-read <secret-path>

SCRIPT_NAME=$(basename "$0")
if [ "$#" != "1" ]; then
    echo usage "$SCRIPT_NAME: <secret-path>"
    exit 1
else
    path=$1
fi
export VAULT_ADDR=${VAULT_ADDR:-http://localhost:8200}

# Get the kv secret from path in json
data=$(vault kv get -format=json $path)
if [ "$?" != "0" ]; then
    exit 1
fi

f=$(echo "$data" | jq -r '.data.format//.data.data.format' 2> /dev/null)
v=$(echo "$data" | jq -r '.data.value//.data.data.value')

# if value format is base64, decode it
if [ "base64" == "$f" ]
then
    echo -n "$v" | base64 --decode
else
    echo -n "$v"
fi
