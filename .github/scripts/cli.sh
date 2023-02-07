#!/bin/bash

echo "Downloading rosetta-cli..."
curl -sSfL https://raw.githubusercontent.com/coinbase/rosetta-cli/master/scripts/install.sh | sh -s
echo "rosetta-cli downloaded"
./bin/rosetta-cli version

echo "Copying rosetta-cli to root..."
mv ./bin/rosetta-cli ./
chmod +x ./rosetta-cli
source ./
rosetta-cli version
echo "Rosetta cli location: $(which rosetta-cli)"
