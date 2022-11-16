# op-rosetta

`op-rosetta` is a Rosetta-compatible API for Optimism Bedrock. It is based on github.com/coinbase/rosetta-geth-sdk.

### Running

```
go install github.com/mdehoog/op-rosetta
export CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'
MODE=ONLINE PORT=5000 NETWORK=mainnet GETH=https://mainnet.optimism.io op-rosetta
```
