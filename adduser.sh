#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: username directory cert_path fetch_address"
    exit 1
}

# Function to generate TLS certificate hash
generate_cert_hash() {
    cert_hash=$(openssl x509 -in "$1" -noout -sha256 -fingerprint | cut -d "=" -f 2 | tr -d ":" | tr '[:upper:]' '[:lower:]')
    echo "$cert_hash"
}

# Function to add a new user
add_user() {
    # Ask for username, and directory path
    username="$1"
    directory_input="$2"
    directory=${directory_input/#\~/${HOME}}  # Replace ~ with $HOME
    echo "Username: $1"
    echo "Directory: $directory"
    echo "Cert file: $3"
    echo "Fetch address: $4"

    # Check if users directory exists
    if [ ! -d "$directory/users" ]; then
        echo "Error: 'users' directory not found."
        exit 1
    fi

    # Set up directory structure for the new user
    user_dir="$directory/users/$username"
    echo "User directory: $user_dir"
    mkdir -p "$user_dir/STORE" "$user_dir/tags" "$user_dir/tags/Archive" "$user_dir/tags/Drafts" "$user_dir/tags/Inbox" "$user_dir/tags/Sent" "$user_dir/tags/Trash"

    # Copy certificate file to cert.pem
    cert_path_input="$3"
    cert_path=${cert_path_input/#\~/${HOME}}  # Replace ~ with $HOME
    cp "$cert_path" "$user_dir/cert.pem"

    # Generate TLS certificate hash
    cert_hash=$(generate_cert_hash "$cert_path")

    # Ask for configuration values
    gemini_fetch_address="$4"
    auto_fetch=false
    timezone="Factory"

    # Create config.yaml file
    cat <<EOF > "$user_dir/config.yaml"
GeminiFetchAddress: $gemini_fetch_address
CertFilename: $user_dir/cert.pem
AutoFetch: $auto_fetch
Timezone: $timezone
# Add other configuration fields as needed
EOF

    # Update users.yaml
    if [ ! -f "$directory/users.yaml" ]; then
        touch $directory/users.yaml
    fi
    echo "$cert_hash: $username" >> $directory/users.yaml

    echo "User $username setup complete."
}

if [ $# == 4 ]; then
    add_user "$@"
else
    usage
fi
