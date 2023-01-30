// Package handlers handles translating transactions into [RosettaTypes.Operation]s.
package handlers

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	// MintOpType is a [RosettaTypes.Operation] type for an Optimism Deposit or "mint" transaction.
	MintOpType = "MINT"
	// BurnOpType is a [RosettaTypes.Operation] type for an Optimism Withdrawal or "burn" transaction.
	BurnOpType = "BURN"

	// An erroneous STOP Type not defined in rosetta-geth-sdk
	StopOpType = "STOP"
)

// Optimism Predeploy Addresses
// See [PredeployedContracts] for more information.
//
// [PredeployedContracts]: https://github.com/ethereum-optimism/optimism/blob/d8e328ae936c6a5f3987c04cbde7bd94403a96a0/specs/predeploys.md
var (
	// The BaseFeeVault predeploy receives the basefees on L2.
	// The basefee is not burnt on L2 like it is on L1.
	// Once the contract has received a certain amount of fees,
	// the ETH can be permissionlessly withdrawn to an immutable address on L1.
	BaseFeeVault = common.HexToAddress("0x4200000000000000000000000000000000000019")

	// The L1FeeVault predeploy receives the L1 portion of the transaction fees.
	// Once the contract has received a certain amount of fees,
	// the ETH can be permissionlessly withdrawn to an immutable address on L1.
	L1FeeVault = common.HexToAddress("0x420000000000000000000000000000000000001a")

	// The L2ToL1MessagePasser stores commitments to withdrawal transactions.
	// When a user is submitting the withdrawing transaction on L1,
	// they provide a proof that the transaction that they withdrew on L2 is in
	// the sentMessages mapping of this contract.
	//
	// Any withdrawn ETH accumulates into this contract on L2 and can be
	// permissionlessly removed from the L2 supply by calling the burn() function.
	L2ToL1MessagePasser = common.HexToAddress("0x4200000000000000000000000000000000000016")
)
