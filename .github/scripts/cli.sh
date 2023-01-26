#!/bin/bash

# downloading cli
echo "Downloading rosetta-cli..."
curl -sSfL https://raw.githubusercontent.com/coinbase/rosetta-cli/master/scripts/install.sh | sh -s
echo "rosetta-cli downloaded"

echo "Copying rosetta-cli to root..."
cp ./bin/rosetta-cli ./
echo "Rosetta cli location: $(which rosetta-cli)"

# echo "start check:data"
# ./bin/rosetta-cli --configuration-file examples/ethereum/rosetta-cli-conf/devnet/config.json check:data --start-block 0
