# Load in the .env file
ifneq ("$(wildcard $(.env))","")
    include .env
endif

# General Config
PWD=$(shell pwd)
NOFILE=100000

# Test all packages
GO_PACKAGES=./...
TEST_SCRIPT=go test ${GO_PACKAGES}
LINT_CONFIG=.golangci.yml
GOIMPORTS_INSTALL=go install golang.org/x/tools/cmd/goimports@latest
GOIMPORTS_CMD=goimports
LINT_SETTINGS=golint,misspell,gocyclo,gocritic,whitespace,goconst,gocognit,bodyclose,unconvert,lll,unparam
GOLINES_INSTALL=go install github.com/segmentio/golines@latest
GOLINES_CMD=golines

.PHONY: clean format lint build run test

# Clean up Rosetta Server binary
clean:
	rm -rf bin/op-rosetta

# Formatting with gofmt
format:
	gofmt -s -w -l .

# Run the golangci-lint linter
lint:
	golangci-lint run -E asciicheck,goimports,misspell ./...

# Build Rosetta Server binary
build:
	env GO111MODULE=on go build -o bin/op-rosetta ./cmd

# Run Rosetta Server
run:
	MODE=${MODE} \
	PORT=${PORT} \
	BLOCKCHAIN=${BLOCKCHAIN} \
	NETWORK=${NETWORK} \
	FILTER=${FILTER} \
	GENESIS_BLOCK_HASH=${GENESIS_BLOCK_HASH} \
	GETH=${GETH} \
	CHAIN_CONFIG=${CHAIN_CONFIG} \
	bin/op-rosetta

# Run tests
test:
	go test -v ./...

