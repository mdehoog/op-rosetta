// Package test contains larger-scoped tests surrounding op-rosetta and it's dependencies.
package test

import (
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
)

var (
	OPCurrency = &RosettaTypes.Currency{
		Symbol:   "OP",
		Decimals: 18,
	}
)
