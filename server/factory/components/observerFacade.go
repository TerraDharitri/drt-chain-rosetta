package components

import (
	"github.com/TerraDharitri/drt-go-chain-proxy/facade"
	"github.com/TerraDharitri/drt-go-chain-proxy/process"
)

// ObserverFacade holds (embeds) several components implemented in proxy-go
type ObserverFacade struct {
	process.Processor
	facade.TransactionProcessor
	facade.BlockProcessor
}

// ComputeShardId computes the shard ID for a given public key
func (facade *ObserverFacade) ComputeShardId(pubKey []byte) uint32 {
	return facade.GetShardCoordinator().ComputeId(pubKey)
}
