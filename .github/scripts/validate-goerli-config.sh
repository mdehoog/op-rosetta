#!/bin/bash

echo "Downloading rosetta-cli..."
curl -sSfL https://raw.githubusercontent.com/coinbase/rosetta-cli/master/scripts/install.sh | sh -s
echo "rosetta-cli downloaded"
./bin/rosetta-cli version

./bin/rosetta-cli configuration:validate configs/optimism/goerli.json
