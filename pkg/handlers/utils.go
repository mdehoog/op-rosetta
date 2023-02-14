package handlers

import (
	"strings"

	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	proxyContractPrefix          = "0x42"
	proxyContractFilter          = "0x420000000000000000000000000000000000"
	implementationContractFilter = "0xc0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3c0d3"
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

func isIdenticalContractAddress(from string, to string) bool {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	proxyContractIndex := from[len(proxyContractFilter):]
	implementationContractIndex := to[len(implementationContractFilter):]
	if strings.Contains(from, proxyContractFilter) && strings.Contains(to, implementationContractFilter) && proxyContractIndex == implementationContractIndex {
		return true
	}

	return false
}
