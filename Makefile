NAME          := IGopher
FILES         := $(wildcard */*.go)
VERSION       := $(shell git describe --always)
BIN_DIR 	  := bin/
BUNDLE_DIR 	  := cmd/igopher/gui-bundler/output/
VUE_DIST_DIR  := resources/static/vue-igopher/dist/

export GO111MODULE=on

## setup: Install required libraries/tools for build tasks
.PHONY: setup
setup:
	@command -v goimports 2>&1 >/dev/null || GO111MODULE=off go get -u -v golang.org/x/tools/cmd/goimports
	@command -v golangci-lint 2>&1 >/dev/null || GO111MODULE=off go get -v github.com/golangci/golangci-lint/cmd/golangci-lint

## fmt: Format all sources files
.PHONY: fmt
fmt: setup 
	goimports -w $(FILES)

## lint: Run all lint related tests against the codebase (will use the .golangci.yml config)
.PHONY: lint
lint: setup
	golangci-lint run

## test: Run the tests against the codebase
.PHONY: test
test:
	go test -v -race ./...

## build: Build the binary for Linux environement
.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 \
		go build \
		-o ./bin/IGopherTUI-linux-amd64 ./cmd/igopher/tui

## build-all: Build binaries for all supported platforms
.PHONY: build-all
build-all:
	env GOOS=linux GOARCH=amd64 \
		go build \
		-o ./bin/IGopherTUI-linux-amd64 \
		./cmd/igopher/tui

	env GOOS=windows GOARCH=amd64 \
		go build \
		-o ./bin/IGopherTUI-windows-amd64.exe \
		./cmd/igopher/tui

	env GOOS=darwin GOARCH=amd64 \
		go build \
		-o ./bin/IGopherTUI-macOS-amd64 \
		./cmd/igopher/tui

## build-vue: Build VueJS project
.PHONY: build-vue
build-vue:
	@if [ @command -v npm 2>&1 >/dev/null ]; then \
		@echo "Npm not found, install NodeJS and retry."; \
		return; \
	fi
	cd ./resources/static/vue-igopher && \
		npm install && \
		npm run build

## bundle: Create astilectron bundle for all supported platforms with embedded ressources
.PHONY: bundle
bundle: build-vue install
	go get github.com/asticode/go-astilectron-bundler/...
	go install github.com/asticode/go-astilectron-bundler/astilectron-bundler
	cd ./cmd/igopher/gui-bundler && \
		mv bind.go bind.go.tmp && \
		astilectron-bundler -c bundler.json && \
		rm bind_*.go windows.syso && \
		mv bind.go.tmp bind.go
	@echo "Done. Executables are located in 'cmd/igopher/gui-bundler/output/' folder"

## release: Build binaries for all platforms for both GUI and TUI
.PHONY: release
release: build-all bundle

## install: Install go dependencies
.PHONY: install
install:
	go get ./...

# vendor: Vendor go modules
.PHONY: vendor
vendor:
	go mod vendor

## coverage: Generates coverage report
.PHONY: coverage
coverage:
	rm -f coverage.out
	go test -v ./... -coverpkg=./... -coverprofile=coverage.out

## clean: Remove binaries (go binaries, bundles and vue dist folder) if they exist
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(BUNDLE_DIR)
	rm -rf $(VUE_DIST_DIR)

.PHONY: all
all: lint test build-vue build-all bundle

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
