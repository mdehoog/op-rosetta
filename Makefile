.PHONY: build
build:
	go build -o bin/op-rosetta .

run-basenet:
	MODE=ONLINE PORT=5000 NETWORK=coinnet \
	GENESIS_BLOCK_HASH=0xc40fbad01e1c6a108a6a07c9a7f1d1051344169915a3a2c40b414bec50736a84 \
	GETH=https://basenet.cbhq.net CHAIN_CONFIG='{ "chainId": 3222583904, "terminalTotalDifficultyPassed": true }' \
	go run .
