package client

import (
	"context"
	"encoding/json"
	"math/big"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	configuration "github.com/coinbase/rosetta-geth-sdk/configuration"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	Eth "github.com/ethereum/go-ethereum"
	EthCommon "github.com/ethereum/go-ethereum/common"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
	EthRpc "github.com/ethereum/go-ethereum/rpc"
)

// InternalClient contains the methods that must be implemented by the parameter to [NewOpClient].
type InternalClient interface {
	// CallContext performs a JSON-RPC call with the given arguments. If the context is
	// canceled before the call has successfully returned, CallContext returns immediately.
	//
	// The result must be a pointer so that package json can unmarshal into it. You
	// can also pass nil, in which case the result is ignored.
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error

	// BatchCallContext sends all given requests as a single batch and waits for the server
	// to return a response for all of them. The wait duration is bounded by the
	// context's deadline.
	//
	// In contrast to CallContext, BatchCallContext only returns errors that have occurred
	// while sending the request. Any error specific to a request is reported through the
	// Error field of the corresponding BatchElem.
	//
	// Note that batch calls may not be executed atomically on the server side.
	BatchCallContext(ctx context.Context, b []EthRpc.BatchElem) error

	// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
	// the current pending state of the backend blockchain. There is no guarantee that this is
	// the true gas limit requirement as other transactions may be added or removed by miners,
	// but it should provide a basis for setting a reasonable default.
	EstimateGas(ctx context.Context, msg Eth.CallMsg) (uint64, error)

	// TODO: remove the methods below that aren't needed

	// Balance returns the balance of an account at a given block.
	Balance(ctx context.Context, account *RosettaTypes.AccountIdentifier, block *RosettaTypes.PartialBlockIdentifier, currencies []*RosettaTypes.Currency) (*RosettaTypes.AccountBalanceResponse, error)
	// BlockAuthor returns the author of a block.
	BlockAuthor(ctx context.Context, blockIndex int64) (string, error)
	// BlockRewardTransaction returns the transaction that rewards the miner of a block.
	BlockRewardTransaction(blockIdentifier *RosettaTypes.BlockIdentifier, miner string, uncles []*EthTypes.Header) *RosettaTypes.Transaction
	// GetBlockReceipts returns [evmClient.RosettaTxReceipt] for a given blockhash.
	GetBlockReceipts(ctx context.Context, blockHash EthCommon.Hash, txs []evmClient.RPCTransaction, baseFee *big.Int) ([]*evmClient.RosettaTxReceipt, error)
	// GetClient returns the internal [evmClient.SDKClient].
	GetClient() *evmClient.SDKClient
	// GetContractCallGasLimit returns the gas limit for a contract call.
	GetContractCallGasLimit(ctx context.Context, toAddress string, fromAddress string, data []byte) (uint64, error)
	// GetContractCurrency returns the currency for a contract.
	GetContractCurrency(addr EthCommon.Address, erc20 bool) (*evmClient.ContractCurrency, error)
	// GetErc20TransferGasLimit returns the gas limit for an ERC20 transfer.
	GetErc20TransferGasLimit(ctx context.Context, toAddress string, fromAddress string, value *big.Int, currency *RosettaTypes.Currency) (uint64, error)
	// GetGasPrice returns the gas price for a transaction.
	GetGasPrice(ctx context.Context, input evmClient.Options) (*big.Int, error)
	// GetLoadedTransaction returns a [evmClient.LoadedTransaction] for a given transaction hash.
	GetLoadedTransaction(ctx context.Context, request *RosettaTypes.BlockTransactionRequest) (*evmClient.LoadedTransaction, error)
	// GetNativeTransferGasLimit returns the gas limit for a native transfer.
	GetNativeTransferGasLimit(ctx context.Context, toAddress string, fromAddress string, value *big.Int) (uint64, error)
	// GetNonce returns the nonce for an account.
	GetNonce(ctx context.Context, input evmClient.Options) (uint64, error)
	// GetRosettaConfig returns the [configuration.RosettaConfig] for the client.
	GetRosettaConfig() configuration.RosettaConfig
	// GetTransactionReceipt returns the [evmClient.RosettaTxReceipt] for a given transaction hash.
	GetTransactionReceipt(ctx context.Context, tx *evmClient.LoadedTransaction) (*evmClient.RosettaTxReceipt, error)
	// GetUncles returns the uncles for a given block.
	GetUncles(ctx context.Context, head *EthTypes.Header, body *evmClient.RPCBlock) ([]*EthTypes.Header, error)
	// ParseOps returns the [RosettaTypes.Operation] for a given [evmClient.LoadedTransaction].
	ParseOps(tx *evmClient.LoadedTransaction) ([]*RosettaTypes.Operation, error)
	// PopulateCrossChainTransactions returns the [RosettaTypes.Transaction] for a given [EthTypes.Block] and [evmClient.LoadedTransaction].
	PopulateCrossChainTransactions(*EthTypes.Block, []*evmClient.LoadedTransaction) ([]*RosettaTypes.Transaction, error)
	// Status returns the current client status.
	Status(ctx context.Context) (*RosettaTypes.BlockIdentifier, int64, *RosettaTypes.SyncStatus, []*RosettaTypes.Peer, error)
	// Submit submits a signed transaction.
	Submit(ctx context.Context, signedTx *EthTypes.Transaction) error
	// TraceBlockByHash returns the [evmClient.FlatCall] for a given block hash.
	TraceBlockByHash(ctx context.Context, blockHash EthCommon.Hash, txs []evmClient.RPCTransaction) (map[string][]*evmClient.FlatCall, error)
	// TraceReplayBlockTransactions returns the [evmClient.FlatCall] for a given block hash.
	TraceReplayBlockTransactions(ctx context.Context, hsh string) (map[string][]*evmClient.FlatCall, error)
	// TraceReplayTransaction returns the [evmClient.FlatCall] for a given transaction hash.
	TraceReplayTransaction(ctx context.Context, hsh string) (json.RawMessage, []*evmClient.FlatCall, error)
	// TraceTransaction constructs trace enriched [emvClient.FlatCall]s for a given [EthCommon.Hash].
	TraceTransaction(ctx context.Context, hash EthCommon.Hash) (json.RawMessage, []*evmClient.FlatCall, error)
}
