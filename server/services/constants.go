package services

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/TerraDharitri/drt-go-chain-core/core"
)

var (
	transactionVersion                                    = 1
	transactionProcessingTypeRelayedV1                    = "RelayedTx"
	transactionProcessingTypeBuiltInFunctionCall          = "BuiltInFunctionCall"
	transactionProcessingTypeMoveBalance                  = "MoveBalance"
	transactionProcessingTypeContractInvoking             = "SCInvoking"
	transactionProcessingTypeContractDeployment           = "SCDeployment"
	amountZero                                            = "0"
	builtInFunctionClaimDeveloperRewards                  = core.BuiltInFunctionClaimDeveloperRewards
	builtInFunctionDCDTTransfer                           = core.BuiltInFunctionDCDTTransfer
	refundGasMessage                                      = "refundedGas"
	argumentsSeparator                                    = "@"
	sendingValueToNonPayableContractDataPrefix            = argumentsSeparator + hex.EncodeToString([]byte("sending value to non payable contract"))
	emptyHash                                             = strings.Repeat("0", 64)
	nodeVersionForOfflineRosetta                          = "N / A"
	systemContractDeployAddress                           = "drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z"
	nativeAsDCDTIdentifier                                = "REWA-000000"
	durationAlarmThresholdBlockServiceGetBlock            = time.Duration(500) * time.Millisecond
	durationAlarmThresholdAccountServiceGetAccountBalance = time.Duration(500) * time.Millisecond
)

const (
	transactionEventSignalError                             = core.SignalErrorOperation
	transactionEventSCDeploy                                = core.SCDeployIdentifier
	transactionEventTransferValueOnly                       = "transferValueOnly"
	transactionEventDCDTTransfer                            = "DCDTTransfer"
	transactionEventDCDTNFTTransfer                         = "DCDTNFTTransfer"
	transactionEventDCDTNFTCreate                           = "DCDTNFTCreate"
	transactionEventDCDTNFTBurn                             = "DCDTNFTBurn"
	transactionEventDCDTNFTAddQuantity                      = "DCDTNFTAddQuantity"
	transactionEventMultiDCDTNFTTransfer                    = "MultiDCDTNFTTransfer"
	transactionEventDCDTLocalBurn                           = core.BuiltInFunctionDCDTLocalBurn
	transactionEventDCDTLocalMint                           = core.BuiltInFunctionDCDTLocalMint
	transactionEventDCDTWipe                                = core.BuiltInFunctionDCDTWipe
	transactionEventClaimDeveloperRewards                   = core.BuiltInFunctionClaimDeveloperRewards
	transactionEventTopicInvalidMetaTransaction             = "meta transaction is invalid"
	transactionEventTopicInvalidMetaTransactionNotEnoughGas = "meta transaction is invalid: not enough gas"

	transactionEventDataExecuteOnDestContext = "ExecuteOnDestContext"
	transactionEventDataAsyncCall            = "AsyncCall"
	transactionEventDataTransferAndExecute   = "TransferAndExecute"
)

const (
	numTopicsOfEventDCDTTransfer                    = 4
	numTopicsPerTransferOfEventMultiDCDTNFTTransfer = 3
	numTopicsOfEventDCDTLocalBurn                   = 3
	numTopicsOfEventDCDTLocalMint                   = 3
	numTopicsOfEventDCDTWipe                        = 4
	numTopicsOfEventDCDTNFTCreate                   = 4
	numTopicsOfEventDCDTNFTBurn                     = 3
	numTopicsOfEventDCDTNFTAddQuantity              = 3
	numTopicsOfEventSCDeployBeforeSirius            = 2
	numTopicsOfEventClaimDeveloperRewards           = 2
	numTopicsOfEventTransferValueOnlyAfterSirius    = 2
)
