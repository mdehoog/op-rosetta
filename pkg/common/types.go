// Package common contains constants and types used across op-rosetta.
package common

import (
	"github.com/ethereum/go-ethereum/common"
)

// Op Types
const (
	// MintOpType is a [RosettaTypes.Operation] type for an Optimism Deposit or "mint" transaction.
	MintOpType = "MINT"
	// BurnOpType is a [RosettaTypes.Operation] type for an Optimism Withdrawal or "burn" transaction.
	BurnOpType = "BURN"
	// An erroneous STOP Type not defined in rosetta-geth-sdk
	StopOpType = "STOP"
)

// Event Topics
const (
	// - If the token is an native token being *sent* to a non-native chain (in either direction), then you should expect an `ERC20BridgeInitiated` event alongside a `Transfer` event FROM the sender TO the bridge contract.
	// - If the token is a non-native token being *sent* to a native chain (in either direction), then you should expect an `ERC20BridgeInitiated` event alongside a `Burn` event FROM the sender.
	// - If the token is a native token being *received* on a non-native chain (in either direction), then you should expect an `ERC20BridgeFinalized` event alongside a `Mint` event TO the sender.
	// - If the token is a non-native token being *received* on a native chain (in either direction), then you should expect an `ERC20BridgeFinalized` event alongside a `Transfer` event FROM the bridge contract TO the sender.

	// TransferEvent is emitted when an ERC20 token is transferred.
	//
	// TransferEvent is emitted in two bridging scenarios:
	// 1. When a native token is being sent to a non-native chain, from the sender to the bridge contract.
	//    Think: Transferring USDC on Ethereum Mainnet to the Optimism bridge contract,
	//    you will see a Transfer event from the sender (you) to the bridge contract.
	// 2. When a non-native token is being sent to a native chain, from the bridge to the sender contract.
	// 	  Think: "Withdrawing" USDC from Optimism to Ethereum Mainnet. You will see a Transfer event
	// 	  from the bridge contract to you (the sender) once the withdrawal is finalized on Mainnet.
	TransferEvent = "Transfer(address,address,uint256)"

	// ERC20BridgeInitiatedEvent is the topic for the ERC20BridgeInitiated event.
	// It is emitted on the originating chain where a bridge is initiated.
	ERC20BridgeInitiatedEvent = "ERC20BridgeInitiated(address,address,address,address,uint256,bytes)"

	// ERC20BridgeFinalizedEvent is the topic for the ERC20BridgeFinalized event.
	// It is emitted on the destination chain where a bridge is finalized.
	ERC20BridgeFinalizedEvent = "ERC20BridgeFinalized(address,address,address,address,uint256,bytes)"

	// Burn event is emitted when a non-native token is being sent to a native chain.
	// For example, consider Bob bridged 100 Token A (native to Ethereum Mainnet) from Ethereum Mainnet to Optimism.
	// Bob now has 100 Token A on Optimism. He then bridges 100 Token A from Optimism to Ethereum Mainnet.
	// In this case, Bob is bridging a non-native token (Token A) from Optimism to the token's native chain (Ethereum Mainnet).
	// In this case, an ERC20BridgeInitiated event will be emitted alongside a Burn event FROM the sender.
	BurnEvent = "Burn(address,uint256)"

	// Mint event is emitted when a non-native token is being sent to a native chain.
	// For example, consider Bob is bridging 100 Token A (native to Ethereum Mainnet) from Ethereum Mainnet to Optimism.
	// Bob will see a Mint event on Optimism TO the sender (Bob).
	MintEvent = "Mint(address,uint256)"
)

// Optimism Deployed Contracts
var (
	// L1StandardBridge is the Ethereum Mainnet Standard Bridge contract deployment.
	//
	L1StandardBridge = common.HexToAddress("0x25ace71c97B33Cc4729CF772ae268934F7ab5fA1")
)

// Optimism Predeploy Addresses (represented as 0x-prefixed hex string)
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

	// The L2StandardBridge predeploy is the contract on L2 used to bridge assets to L1.
	L2StandardBridge = common.HexToAddress("0x4200000000000000000000000000000000000010")
)
