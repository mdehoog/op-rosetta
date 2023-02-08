package handlers

import (
	"github.com/coinbase/rosetta-sdk-go/types"
)

func GenerateOp(opIndex int64, relatedOps []*types.OperationIdentifier, opType string, opStatus string, address string, amount *types.Amount, metadata map[string]interface{}) *types.Operation {
	return &types.Operation{
		OperationIdentifier: &types.OperationIdentifier{
			Index: opIndex,
		},
		RelatedOperations: relatedOps,
		Type:              opType,
		Status:            types.String(opStatus),
		Account: &types.AccountIdentifier{
			Address: address,
		},
		Amount:   amount,
		Metadata: metadata,
	}
}
