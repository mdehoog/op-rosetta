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
		from := evmClient.MustChecksum(call.From.String())
		to := evmClient.MustChecksum(call.To.String())
		value := call.Value
		metadata := map[string]interface{}{}

		// Handle the case where not all operation statuses are successful
		opStatus := sdkTypes.SuccessStatus
		if call.Revert {
			opStatus = sdkTypes.FailureStatus
			metadata["error"] = call.ErrorMessage
		}

		fromOp := &RosettaTypes.Operation{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: int64(len(ops) + startIndex),
			},
			Type:   opType,
			Status: RosettaTypes.String(opStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: from,
			},
			Amount:   evmClient.Amount(new(big.Int).Neg(value), sdkTypes.Currency),
			Metadata: metadata,
		}
		if _, ok := destroyedAccountBalance[from]; ok && opStatus == sdkTypes.SuccessStatus {
			destroyedAccountBalance[from] = new(big.Int).Sub(destroyedAccountBalance[from], value)
		}
		ops = append(ops, fromOp)

		// Add to the destroyed account balance if SELFDESTRUCT, and overwrite existing balance.
		if opType == sdkTypes.SelfDestructOpType {
			destroyedAccountBalance[from] = new(big.Int)

			// If destination of of SELFDESTRUCT is self, we should skip.
			// In the EVM, the balance is reset after the balance is increased on the destination, so this is a no-op.
			if from == to {
				continue
			}
		}

		// If the account is resurrected, we remove it from the destroyed account balance map.
		if sdkTypes.CreateType(opType) {
			delete(destroyedAccountBalance, to)
		}

		lastOpIndex := ops[len(ops)-1].OperationIdentifier.Index
		toOp := &RosettaTypes.Operation{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: lastOpIndex + 1,
			},
			RelatedOperations: []*RosettaTypes.OperationIdentifier{
				{
					Index: lastOpIndex,
				},
			},
			Type:   opType,
			Status: RosettaTypes.String(opStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: to,
			},
			Amount:   evmClient.Amount(new(big.Int).Abs(value), sdkTypes.Currency),
			Metadata: metadata,
		}
		if _, ok := destroyedAccountBalance[to]; ok && opStatus == sdkTypes.SuccessStatus {
			destroyedAccountBalance[to] = new(big.Int).Add(destroyedAccountBalance[to], value)
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

		ops = append(ops, &RosettaTypes.Operation{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: ops[len(ops)-1].OperationIdentifier.Index + 1,
			},
			Type:   sdkTypes.DestructOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: acct,
			},
			Amount: evmClient.Amount(new(big.Int).Neg(balance), sdkTypes.Currency),
		})
	}

	return ops
}
