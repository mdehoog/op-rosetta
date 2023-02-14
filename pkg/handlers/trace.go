package handlers

import (
	"log"
	"math/big"
	"strings"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
)

func TraceOps(
	calls []*evmClient.FlatCall,
	startIndex int,
) []*RosettaTypes.Operation { // nolint: gocognit
	var ops []*RosettaTypes.Operation
	if len(calls) == 0 {
		return ops
	}

	destroyedAccountBalance := make(map[string]*big.Int)
	for _, call := range calls {
		opType := strings.ToUpper(call.Type)
		fromAddress := evmClient.MustChecksum(call.From.String())
		toAddress := evmClient.MustChecksum(call.To.String())
		if strings.Contains(fromAddress, proxyContractPrefix) && isIdenticalContractAddress(fromAddress, toAddress) {
			toAddress = fromAddress
		}
		value := call.Value
		metadata := map[string]interface{}{}

		// Handle the case where not all operation statuses are successful
		opStatus := sdkTypes.SuccessStatus
		if call.Revert {
			opStatus = sdkTypes.FailureStatus
			metadata["error"] = call.ErrorMessage
		}

		// Generate "from" operation
		fromOpIndex := int64(len(ops) + startIndex)
		fromAmount := evmClient.Amount(new(big.Int).Neg(value), sdkTypes.Currency)
		fromOp := GenerateOp(fromOpIndex, nil, opType, opStatus, fromAddress, fromAmount, metadata)
		if _, ok := destroyedAccountBalance[fromAddress]; ok && opStatus == sdkTypes.SuccessStatus {
			destroyedAccountBalance[fromAddress] = new(big.Int).Sub(destroyedAccountBalance[fromAddress], value)
		}
		ops = append(ops, fromOp)

		// Add to the destroyed account balance if SELFDESTRUCT, and overwrite existing balance.
		if opType == sdkTypes.SelfDestructOpType {
			destroyedAccountBalance[fromAddress] = new(big.Int)

			// If destination of of SELFDESTRUCT is self, we should skip.
			// In the EVM, the balance is reset after the balance is increased on the destination, so this is a no-op.
			if fromAddress == toAddress {
				continue
			}
		}

		// If the account is resurrected, we remove it from the destroyed account balance map.
		if sdkTypes.CreateType(opType) {
			delete(destroyedAccountBalance, toAddress)
		}

		// Generate "to" operation
		lastOpIndex := ops[len(ops)-1].OperationIdentifier.Index
		toOpIndex := lastOpIndex + 1
		toRelatedOps := []*RosettaTypes.OperationIdentifier{
			{
				Index: lastOpIndex,
			},
		}
		toAmount := evmClient.Amount(new(big.Int).Abs(value), sdkTypes.Currency)
		toOp := GenerateOp(toOpIndex, toRelatedOps, opType, opStatus, toAddress, toAmount, metadata)
		if _, ok := destroyedAccountBalance[toAddress]; ok && opStatus == sdkTypes.SuccessStatus {
			destroyedAccountBalance[toAddress] = new(big.Int).Add(destroyedAccountBalance[toAddress], value)
		}
		ops = append(ops, toOp)
	}

	// Zero-out all destroyed accounts that are removed during transaction finalization.
	for acct, balance := range destroyedAccountBalance {
		if _, ok := evmClient.ChecksumAddress(acct); !ok || balance.Sign() == 0 {
			continue
		}

		if balance.Sign() < 0 {
			log.Fatalf("negative balance for suicided account %s: %s\n", acct, balance.String())
		}

		// Generate "destruct" operation
		destructOpIndex := ops[len(ops)-1].OperationIdentifier.Index + 1
		destructOpType := sdkTypes.DestructOpType
		destructOpStatus := RosettaTypes.String(sdkTypes.SuccessStatus)
		address := acct
		amount := evmClient.Amount(new(big.Int).Neg(balance), sdkTypes.Currency)
		destructOp := GenerateOp(destructOpIndex, nil, destructOpType, *destructOpStatus, address, amount, nil)
		ops = append(ops, destructOp)
	}

	return ops
}
