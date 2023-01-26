# Load in the .env file
include .env

# General Config
PWD=$(shell pwd)
NOFILE=100000

all: clean build test lint

clean:
	rm -rf bin/op-rosetta

build:
	env GO111MODULE=on go build -o bin/op-rosetta ./cmd

test:
	go test -v ./...

lint:
	golangci-lint run -E asciicheck,goimports,misspell ./...

.PHONY: \
	test \
	lint \
	build

# Runs an instance of `op-rosetta` configured for Optimism Mainnet
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
	GENESIS_BLOCK_HASH=0xa89b19033c8b43365e244f425a7e4acb5bae21d1893e1be0eb8cddeb29950d72 \
	bin/op-rosetta

# Runs an instance of `op-rosetta` configured for Optimism Goerli
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

