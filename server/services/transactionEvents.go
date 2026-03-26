package services

import (
	"fmt"
	"math/big"
)

type eventTransferValueOnly struct {
	sender   string
	receiver string
	value    string
}

type eventDCDT struct {
	senderAddress   string
	receiverAddress string
	otherAddress    string
	identifier      string
	nonceAsBytes    []byte
	value           string
}

// newEventDCDTFromBasicTopics creates an eventDCDT from the given topics. The following topics are expected:
// - topic 0: the identifier of the token
// - topic 1: the nonce of the token
// - topic 2: the value of the token
func newEventDCDTFromBasicTopics(topics [][]byte) (*eventDCDT, error) {
	if len(topics) < 3 {
		return nil, fmt.Errorf("newEventDCDTFromBasicTopics: bad number of topics: %d", len(topics))
	}

	identifier := topics[0]
	nonceAsBytes := topics[1]
	valueBytes := topics[2]
	value := big.NewInt(0).SetBytes(valueBytes)

	return &eventDCDT{
		identifier:   string(identifier),
		nonceAsBytes: nonceAsBytes,
		value:        value.String(),
	}, nil
}

// getBaseIdentifier returns the token identifier for fungible tokens, and the collection identifier for SFTs, NFTs and MetaDCDTs
func (event *eventDCDT) getBaseIdentifier() string {
	return event.identifier
}

// getExtendedIdentifier returns the "full" token identifier for all types of DCDTs
func (event *eventDCDT) getExtendedIdentifier() string {
	if len(event.nonceAsBytes) > 0 {
		return fmt.Sprintf("%s-%x", event.identifier, event.nonceAsBytes)
	}

	return event.identifier
}

type eventSCDeploy struct {
	contractAddress string
	deployerAddress string
}

type eventClaimDeveloperRewards struct {
	value           string
	receiverAddress string
}
