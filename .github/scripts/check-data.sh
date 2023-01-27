#!/bin/bash

./.github/scripts/cli.sh

ROSETTA_CONFIGURATION_FILE=configs/optimism/goerli.json ./bin/rosetta-cli check:data configs/optimism/goerli.json
