#!/bin/bash
# Script used to build GUI version and bundle executables for Windows/Linux/MacOS

# Check if pwd is located inside the SkillsList folder
if [ $(basename "$PWD") != "gui-bundler" ]; then
    cd "${0%/*}/../cmd/igopher/gui-bundler"
    if [ $(basename "$PWD") != "gui-bundler" ]; then
        echo "Invalid current directory! Please cd to the IGopher directory and re-run this script."
        exit 1
    fi
fi

# Download and install go-astilectron-bundler
go get github.com/asticode/go-astilectron-bundler/...
go install github.com/asticode/go-astilectron-bundler/astilectron-bundler

# Install dependencies 
go get ../../../...

# Rename default bind.go to tmp file 
mv bind.go bind.go.tmp

# Execute astilectron-bundler
astilectron-bundler -c bundler.json

# Delete generated files and restore bind.go
rm bind_*.go windows.syso
mv bind.go.tmp bind.go
echo "Done. Executables are located in 'cmd/igopher/gui-bundler/output/' folder"