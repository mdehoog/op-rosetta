# Load in the .env file
# include .env
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

# Run the full pipeline
all: clean format build test lint
.PHONY: \
	test \
	tests \
	format \
	lint \
	build \
	clean

# Clean up built binaries
clean:
	rm -rf bin/op-rosetta

format:
	gofmt -s -w -l .

# Build the `op-rosetta` binary
build:
	env GO111MODULE=on go build -o bin/op-rosetta ./cmd

# Comprehensive tests
test: tests
tests: unit-tests integration-tests
unit-tests:
	go test -v ./...
integration-tests: config-validation
config-validation: run-optimism-mainnet-validate-config run-optimism-goerli-validate-config

# TODO: Add the `check:data` command to the pipeline
# TODO: Requires node env var configuration in the github repository for the actions to run tests successfully
check-data: run-optimism-mainnet-data-check run-optimism-goerli-data-check

# Run the golangci-lint linter
lint:
	golangci-lint run -E asciicheck,goimports,misspell ./...


##################################################################################
## GOERLI GOERLI GOERLI GOERLI GOERLI GOERLI GOERLI GOERLI GOERLI GOERLI GOERLI ##
##################################################################################

# Runs rosetta-cli configuration:validate against the optimism goerli configuration
run-optimism-goerli-validate-config:
	ROSETTA_CONFIGURATION_FILE=configs/optimism/goerli.json rosetta-cli configuration:validate configs/optimism/goerli.json

# Runs the rosetta-cli check:data command with the optimism goerli configuration
run-optimism-goerli-data-check:
	ROSETTA_CONFIGURATION_FILE=configs/optimism/goerli.json rosetta-cli check:data configs/optimism/goerli.json

# Runs an instance of `op-rosetta` configured for Optimism Goerli
# For the genesis block hash, see:
# https://github.com/ethereum-optimism/optimism/blob/5e8bc3d5b4f36f0192b22b032e25b09f23cd0985/op-node/chaincfg/chains.go#L49
run-optimism-goerli:
	CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'	\
	MODE=ONLINE \
	PORT=8080 \
	BLOCKCHAIN=Optimism \
	NETWORK=Goerli \
	GETH=${OPTIMISM_GOERLI_NODE} \
	ENABLE_TRACE_CACHE=true \
    ENABLE_GETH_TRACER=true \
	GENESIS_BLOCK_HASH=0x0f783549ea4313b784eadd9b8e8a69913b368b7366363ea814d7707ac505175f \
	bin/op-rosetta

#####################################################################################
## MAINNET MAINNET MAINNET MAINNET MAINNET MAINNET MAINNET MAINNET MAINNET MAINNET ##
#####################################################################################

# Runs rosetta-cli configuration:validate against the optimism mainnet configuration
run-optimism-mainnet-validate-config:
	ROSETTA_CONFIGURATION_FILE=configs/optimism/mainnet.json rosetta-cli configuration:validate configs/optimism/mainnet.json

# Runs the rosetta-cli check:data command with the optimism mainnet configuration
run-optimism-mainnet-data-check:
	ROSETTA_CONFIGURATION_FILE=configs/optimism/mainnet.json rosetta-cli check:data configs/optimism/mainnet.json


# TODO: Set the GENESIS_BLOCK_HASH value for this command
# Runs an instance of `op-rosetta` configured for Optimism Mainnet
# For the genesis block hash, see:
# https://github.com/ethereum-optimism/optimism/blob/5e8bc3d5b4f36f0192b22b032e25b09f23cd0985/op-node/chaincfg/chains.go
run-optimism-mainnet:
	make build \
	CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'	\
	MODE=ONLINE \
	PORT=8080 \
	BLOCKCHAIN=Optimism \
	NETWORK=Mainnet \
	ENABLE_TRACE_CACHE=true \
    ENABLE_GETH_TRACER=true \
	GETH=${OPTIMISM_MAINNET_NODE} \
	GENESIS_BLOCK_HASH= \
	bin/op-rosetta

