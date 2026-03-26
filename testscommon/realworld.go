package testscommon

import (
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	hasherFactory "github.com/TerraDharitri/drt-go-chain-core/hashing/factory"
	marshalFactory "github.com/TerraDharitri/drt-go-chain-core/marshal/factory"
)

// RealWorldBech32PubkeyConverter is a bech32 converter, to be used in tests
var RealWorldBech32PubkeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "drt")

// RealWorldBlake2bHasher is a blake2b hasher, to be used in tests
var RealWorldBlake2bHasher, _ = hasherFactory.NewHasher("blake2b")

// MarshalizerForHashing is a gogo protobuf marshalizer, to be used in tests
var MarshalizerForHashing, _ = marshalFactory.NewMarshalizer("gogo protobuf")
