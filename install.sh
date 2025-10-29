#!/bin/bash

# Password Generator CLI Installation Script

set -e

REPO_URL="https://github.com/thetanav/passwordgen.git"
TEMP_DIR=$(mktemp -d)

echo "Installing passwordgen cli from thetanav CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    echo "Visit https://golang.org/dl/ for installation instructions."
    exit 1
fi

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "Error: Git is not installed. Please install Git first."
    exit 1
fi

# Clone the repository
echo "Downloading source code..."
git clone "$REPO_URL" "$TEMP_DIR"

# Change to the cloned directory
cd "$TEMP_DIR"

# Build the binary
echo "Building the application..."
go build -o passwordgen .

# Determine installation directory
if [[ -w /usr/local/bin ]]; then
    INSTALL_DIR="/usr/local/bin"
elif [[ -d "$HOME/bin" ]]; then
    INSTALL_DIR="$HOME/bin"
else
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Install the binary
echo "Installing to $INSTALL_DIR..."
mv passwordgen "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/passwordgen"

# Clean up
cd /
rm -rf "$TEMP_DIR"

# Check if INSTALL_DIR is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Note: $INSTALL_DIR is not in your PATH."
    echo "Add the following line to your shell profile (.bashrc, .zshrc, etc.):"
    echo "export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "Installation complete! Run 'passwordgen' to start the application."