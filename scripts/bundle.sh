#!/bin/bash
# Script used to build GUI version and bundle executables for Windows/Linux/MacOS
# Require Go and Node/Npm installed
# You must locate your terminal into the scripts sub-folder before running this script

# Check if pwd is located inside the "scripts" sub-folder
if [ $(basename "$PWD") != "scripts" ]; then
    cd "scripts"
    if [ $(basename "$PWD") != "scripts" ]; then
        echo "Invalid current directory! Please cd to the IGopher scripts sub-directory and re-run this script."
        exit 1
    fi
fi

cd "../resources/static/vue-igopher"

# Install node dependencies and build vue-igopher
npm install
npm run build

cd "../../../cmd/igopher/gui-bundler"

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