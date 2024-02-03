#!/bin/bash

# Prompt the user to enter the hostname
read -p "What directory should the CGI executables be built to? Enter path:" hostname

# Iterate over each Go file and build with the provided hostname
for file in *.go; do
    if [[ -f "$file" ]]; then
        name="${file%.go}"  # Extract the filename without the .go extension
        # Build the binary with the provided hostname
        go build -o "$hostname/$name" "$file"
    fi
done
