package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	hasherFactory "github.com/TerraDharitri/drt-go-chain-core/hashing/factory"
	marshalFactory "github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
	"github.com/TerraDharitri/drt-go-chain/sharding"
	"github.com/TerraDharitri/drt-go-chain-proxy/data"
	"github.com/TerraDharitri/drt-go-chain-proxy/observer"
	"github.com/TerraDharitri/drt-go-chain-proxy/process"
	processFactory "github.com/TerraDharitri/drt-go-chain-proxy/process/factory"
	"github.com/TerraDharitri/drt-chain-rosetta/server/factory/components"
	"github.com/TerraDharitri/drt-chain-rosetta/server/provider"
	"github.com/TerraDharitri/drt-chain-rosetta/server/resources"
)

const (
	hasherType                = "blake2b"
	marshalizerForHashingType = "gogo protobuf"
	pubKeyLength              = 32
	bech32Prefix              = "drt"

	notApplicableConfigurationFilePath   = "not applicable"
	notApplicableFullHistoryNodesMessage = "not applicable"

	requestTimeoutInSeconds = 60
)

type ArgsCreateNetworkProvider struct {
	IsOffline                   bool
	NumShards                   uint32
	ObservedActualShard         uint32
	ObservedProjectedShard      uint32
	ObservedProjectedShardIsSet bool
	ObserverUrl                 string
	BlockchainName              string
	NetworkID                   string
	NetworkName                 string
	GasPerDataByte              uint64
	GasPriceModifier            float64
	GasLimitCustomTransfer      uint64
	MinGasPrice                 uint64
	MinGasLimit                 uint64
	ExtraGasLimitGuardedTx      uint64
	ExtraGasLimitRelayedTxV3    uint64
	NativeCurrencySymbol        string
	CustomCurrencies            []resources.Currency
	GenesisBlockHash            string
	GenesisTimestamp            int64
	FirstHistoricalEpoch        uint32
	NumHistoricalEpochs         uint32
	ShouldHandleContracts       bool
}

// CreateNetworkProvider creates a network provider
func CreateNetworkProvider(args ArgsCreateNetworkProvider) (NetworkProvider, error) {
	shardCoordinator, err := sharding.NewMultiShardCoordinator(args.NumShards, args.ObservedActualShard)
	if err != nil {
		return nil, err
	}

	pubKeyConverter, err := pubkeyConverter.NewBech32PubkeyConverter(pubKeyLength, bech32Prefix)
	if err != nil {
		return nil, err
	}

	observers := []*data.NodeData{
		{
			ShardId:  args.ObservedActualShard,
			Address:  args.ObserverUrl,
			IsSynced: true,
		},
	}

	observersProvider, err := observer.NewSimpleNodesProvider(
		observers,
		notApplicableConfigurationFilePath,
		shardCoordinator.NumberOfShards(),
	)
	if err != nil {
		return nil, err
	}

	disabledObserversProvider := observer.NewDisabledNodesProvider(notApplicableFullHistoryNodesMessage)
	if err != nil {
		return nil, err
	}

	baseProcessor, err := process.NewBaseProcessor(
		requestTimeoutInSeconds,
		shardCoordinator,
		observersProvider,
		disabledObserversProvider,
		pubKeyConverter,
		true,
	)
	if err != nil {
		return nil, err
	}

	hasher, err := hasherFactory.NewHasher(hasherType)
	if err != nil {
		return nil, err
	}

	marshalizerForHashing, err := marshalFactory.NewMarshalizer(marshalizerForHashingType)
	if err != nil {
		return nil, err
	}

	transactionProcessor, err := processFactory.CreateTransactionProcessor(
		baseProcessor,
		pubKeyConverter,
		hasher,
		marshalizerForHashing,
		false,
	)
	if err != nil {
		return nil, err
	}

	blockProcessor, err := process.NewBlockProcessor(baseProcessor)
	if err != nil {
		return nil, err
	}

	return provider.NewNetworkProvider(provider.ArgsNewNetworkProvider{
		IsOffline:                   args.IsOffline,
		ObservedActualShard:         args.ObservedActualShard,
		ObservedProjectedShard:      args.ObservedProjectedShard,
		ObservedProjectedShardIsSet: args.ObservedProjectedShardIsSet,
		ObserverUrl:                 args.ObserverUrl,
		BlockchainName:              args.BlockchainName,
		NetworkID:                   args.NetworkID,
		NetworkName:                 args.NetworkName,
		GasPerDataByte:              args.GasPerDataByte,
		GasPriceModifier:            args.GasPriceModifier,
		GasLimitCustomTransfer:      args.GasLimitCustomTransfer,
		MinGasPrice:                 args.MinGasPrice,
		MinGasLimit:                 args.MinGasLimit,
		ExtraGasLimitGuardedTx:      args.ExtraGasLimitGuardedTx,
		ExtraGasLimitRelayedTxV3:    args.ExtraGasLimitRelayedTxV3,
		NativeCurrencySymbol:        args.NativeCurrencySymbol,
		CustomCurrencies:            args.CustomCurrencies,
		GenesisBlockHash:            args.GenesisBlockHash,
		GenesisTimestamp:            args.GenesisTimestamp,
		FirstHistoricalEpoch:        args.FirstHistoricalEpoch,
		NumHistoricalEpochs:         args.NumHistoricalEpochs,
		ShouldHandleContracts:       args.ShouldHandleContracts,

		ObserverFacade: &components.ObserverFacade{
			Processor:            baseProcessor,
			TransactionProcessor: transactionProcessor,
			BlockProcessor:       blockProcessor,
		},

		Hasher:                hasher,
		MarshalizerForHashing: marshalizerForHashing,
		PubKeyConverter:       pubKeyConverter,
	})
}
