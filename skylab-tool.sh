#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 [--init]"
    echo "Options:"
    echo "  --init    Perform initialization (create directories, etc.)"
    echo "  --adduser Add a new user"
    exit 1
}

# Function to set up Skylab initialization
skylab_init() {
    # Create users directory if it doesn't exist
    read -p "Enter directory path for Skylab: " directory_input
    read -p "Enter hostname: " hostname
    directory=${directory_input/#\~/${HOME}}  # Replace ~ with $HOME
    if [ ! -d "$directory/users" ]; then
        mkdir -p "$directory/users"
        mkdir -p "$directory/$hostname"
        touch "$directory/users.yaml"
        touch "$directory/site_config.yaml"
        cat <<EOF > "$directory/site_config.yaml"
ExposeSettings: true
EOF
    fi
    echo "Initialization complete."
}

# Function to generate TLS certificate hash
generate_cert_hash() {
    cert_hash=$(openssl x509 -in "$1" -noout -sha256 -fingerprint | cut -d "=" -f 2 | tr -d ":" | tr '[:upper:]' '[:lower:]')
    echo "$cert_hash"
}

# Function to add a new user
add_user() {
    # Ask for username, and directory path
    read -p "Enter username: " username
    read -p "Enter directory path for Skylab. This directory should already exist and should contain 'users' and 'users.yaml': " directory_input
    directory=${directory_input/#\~/${HOME}}  # Replace ~ with $HOME

    # Check if users directory exists
    if [ ! -d "$directory/users" ]; then
        echo "Error: 'users' directory not found. Run with --init flag first."
        exit 1
    fi

    # Set up directory structure for the new user
    user_dir="$directory/users/$username"
    mkdir -p "$user_dir/STORE" "$user_dir/tags" "$user_dir/tags/Archive" "$user_dir/tags/Drafts" "$user_dir/tags/Inbox" "$user_dir/tags/Sent" "$user_dir/tags/Trash"

    # Copy certificate file to cert.pem
    read -p "Enter path to misfin keypair PEM file (this script will copy it): " cert_path_input
    cert_path=${cert_path_input/#\~/${HOME}}  # Replace ~ with $HOME
    cp "$cert_path" "$user_dir/cert.pem"

    # Generate TLS certificate hash
    cert_hash=$(generate_cert_hash "$cert_path")

    # Ask for configuration values
    read -p "Enter GeminiFetchAddress: " gemini_fetch_address
    read -p "Are you using a GMAP-style gembox? (true/false): " GMAP
    read -p "Enter AutoFetch (true/false): " auto_fetch
    read -p "Enter Timezone: " timezone

    # Create config.yaml file
    cat <<EOF > "$user_dir/config.yaml"
GeminiFetchAddress: $gemini_fetch_address
CertFilename: $user_dir/cert.pem
AutoFetch: $auto_fetch
Timezone: $timezone
GMAP: $GMAP
DisplayNumber: 4
# Add other configuration fields as needed
EOF

    # Update users.yaml
    if [ ! -f "$directory/users.yaml" ]; then
        touch $directory/users.yaml
    fi
    echo "$cert_hash: $username" >> $directory/users.yaml

    echo "User $username setup complete."
}
# Main script
if [ "$1" == "--init" ]; then
    skylab_init
elif [ "$1" == "--adduser" ]; then
    add_user
else
    usage
fi
