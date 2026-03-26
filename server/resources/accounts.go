package resources

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
)

// AccountApiResponse is an API resource
type AccountApiResponse struct {
	resourceApiResponse
	Data AccountOnBlock `json:"data"`
}

// AccountOnBlock defines an account resource
type AccountOnBlock struct {
	Account          Account          `json:"account"`
	BlockCoordinates BlockCoordinates `json:"blockInfo"`
}

// Account defines an account resource
type Account struct {
	Address string `json:"address"`
	Nonce   uint64 `json:"nonce"`
	Balance string `json:"balance"`
}

// AccountDCDTBalanceApiResponse is an API resource
type AccountDCDTBalanceApiResponse struct {
	resourceApiResponse
	Data AccountDCDTBalanceApiResponsePayload `json:"data"`
}

// AccountDCDTBalanceApiResponsePayload is an API resource
type AccountDCDTBalanceApiResponsePayload struct {
	TokenData        AccountDCDTTokenData `json:"tokenData"`
	BlockCoordinates BlockCoordinates     `json:"blockInfo"`
}

// AccountDCDTTokenData is an API resource
type AccountDCDTTokenData struct {
	Identifier string `json:"tokenIdentifier"`
	Balance    string `json:"balance"`
}

// AccountBalanceOnBlock defines an account resource
type AccountBalanceOnBlock struct {
	Balance          string
	Nonce            core.OptionalUint64
	BlockCoordinates BlockCoordinates
}
