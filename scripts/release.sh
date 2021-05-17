#!/bin/bash
# Script used to build GUI and TUI version for Windows/Linux/MacOS

# Check if pwd is located inside IGopher root directory
if [ $(basename "$PWD") != "igopher" ]; then
    cd "${0%/*}/.."
    if [ $(basename "$PWD") != "gui-bundler" ]; then
        echo "Invalid current directory! Please cd to the IGopher directory and re-run this script."
        exit 1
    fi
fi

# Call bundle.sh script to build GUI executables
./scripts/bundle.sh

# Build TUI executables fot all OS
env GOOS=linux GOARCH=amd64 go build -o ./bin/IGopherTUI-linux-amd64 ./cmd/igopher/tui
env GOOS=windows GOARCH=amd64 go build -o ./bin/IGopherTUI-windows-amd64.exe ./cmd/igopher/tui
env GOOS=darwin GOARCH=amd64 go build -o ./bin/IGopherTUI-macOS-amd64 ./cmd/igopher/tui

echo "Done. TUI executables are located in 'bin/' folder"
