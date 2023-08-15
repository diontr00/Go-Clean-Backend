#!/bin/bash

# Read the Makefile line by line and extract target names and comments
while IFS='' read -r line || [[ -n "$line" ]]; do
    summary=$(echo "$line" | grep -Eo '^#.*')
    target=$(echo "$line" | grep -Eo '^[a-zA-Z0-9_-]+:')
    comment=$(echo "$line" | grep -Eo '##.*$' | sed 's/##//')

    if [[ -n "$summary" ]]; then
        echo
        printf "%-20s %s\n" "$summary"
        echo "---"
    fi

    if [[ -n "$target" && -n "$comment" ]]; then
        printf "%-20s %s\n" "$target" "$comment"
    fi
done < Makefile
