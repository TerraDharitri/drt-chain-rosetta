package services

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
)

type transactionEventsController struct {
	provider NetworkProvider
}

func newTransactionEventsController(provider NetworkProvider) *transactionEventsController {
	return &transactionEventsController{
		provider: provider,
	}
}

func (controller *transactionEventsController) extractEventSCDeploy(tx *transaction.ApiTransactionResult) ([]*eventSCDeploy, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventSCDeploy)
	typedEvents := make([]*eventSCDeploy, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics < numTopicsOfEventSCDeployBeforeSirius {
			// Before Sirius, there are 2 topics: contract address, deployer address.
			// After Sirius, there are 3 topics: contract address, deployer address, codehash (not used).
			return nil, fmt.Errorf("%w: bad number of topics for SCdeploy event = %d", errCannotRecognizeEvent, numTopics)
		}

		// "event.Address" is same as "event.Topics[0]"" (the address of the deployed contract).
		contractAddress := event.Address
		deployerPubKey := event.Topics[1]
		deployerAddress := controller.provider.ConvertPubKeyToAddress(deployerPubKey)

		typedEvents = append(typedEvents, &eventSCDeploy{
			contractAddress: contractAddress,
			deployerAddress: deployerAddress,
		})
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventTransferValueOnly(tx *transaction.ApiTransactionResult) ([]*eventTransferValueOnly, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventTransferValueOnly)
	typedEvents := make([]*eventTransferValueOnly, 0)

	for _, event := range rawEvents {
		typedEvent, err := controller.decideEffectiveEventTransferValueOnlyAfterSirius(event)
		if err != nil {
			return nil, err
		}

		if typedEvent != nil {
			typedEvents = append(typedEvents, typedEvent)
		}
	}

	return typedEvents, nil
}

// See: https://github.com/TerraDharitri/drt-specs/blob/main/releases/protocol/release-specs-v1.6.0-Sirius.md#17-logs--events-changes-5490
func (controller *transactionEventsController) decideEffectiveEventTransferValueOnlyAfterSirius(event *transaction.Events) (*eventTransferValueOnly, error) {
	numTopics := len(event.Topics)

	if numTopics != numTopicsOfEventTransferValueOnlyAfterSirius {
		return nil, fmt.Errorf("%w: bad number of topics for 'transferValueOnly' = %d", errCannotRecognizeEvent, numTopics)
	}

	valueBytes := event.Topics[0]
	receiverPubKey := event.Topics[1]

	if len(valueBytes) == 0 {
		return nil, nil
	}

	eventData := string(event.Data)
	if eventData != transactionEventDataExecuteOnDestContext && eventData != transactionEventDataAsyncCall && eventData != transactionEventDataTransferAndExecute {
		// Ineffective event, since the balance change is already captured by a SCR.
		return nil, nil
	}

	sender := event.Address
	senderPubKey, err := controller.provider.ConvertAddressToPubKey(sender)
	if err != nil {
		return nil, err
	}

	isIntrashard := controller.provider.ComputeShardIdOfPubKey(senderPubKey) == controller.provider.ComputeShardIdOfPubKey(receiverPubKey)
	if !isIntrashard {
		// Ineffective event, since the balance change is already captured by a SCR.
		return nil, nil
	}

	receiver := controller.provider.ConvertPubKeyToAddress(receiverPubKey)
	value := big.NewInt(0).SetBytes(valueBytes)

	return &eventTransferValueOnly{
		sender:   sender,
		receiver: receiver,
		value:    value.String(),
	}, nil
}

func (controller *transactionEventsController) hasAnySignalError(tx *transaction.ApiTransactionResult) bool {
	if !controller.hasEvents(tx) {
		return false
	}

	for _, event := range tx.Logs.Events {
		isSignalError := event.Identifier == transactionEventSignalError
		if isSignalError {
			return true
		}
	}

	return false
}

func (controller *transactionEventsController) hasSignalErrorOfSendingValueToNonPayableContract(tx *transaction.ApiTransactionResult) bool {
	if !controller.hasEvents(tx) {
		return false
	}

	for _, event := range tx.Logs.Events {
		isSignalError := event.Identifier == transactionEventSignalError
		dataAsString := string(event.Data)
		dataMatchesError := strings.HasPrefix(dataAsString, sendingValueToNonPayableContractDataPrefix)

		if isSignalError && dataMatchesError {
			return true
		}
	}

	return false
}

func (controller *transactionEventsController) extractEventsDCDTOrDCDTNFTTransfers(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEventsDCDTTransfer := controller.findManyEventsByIdentifier(tx, transactionEventDCDTTransfer)
	rawEventsDCDTNFTTransfer := controller.findManyEventsByIdentifier(tx, transactionEventDCDTNFTTransfer)
	rawEventsMultiDCDTNFTTransfer := controller.findManyEventsByIdentifier(tx, transactionEventMultiDCDTNFTTransfer)

	typedEvents := make([]*eventDCDT, 0)

	// First, handle single transfers
	for _, event := range append(rawEventsDCDTTransfer, rawEventsDCDTNFTTransfer...) {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTTransfer {
			return nil, fmt.Errorf("%w: bad number of topics for (DCDT|DCDTNFT)Transfer event = %d", errCannotRecognizeEvent, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		receiverPubkey := event.Topics[3]
		typedEvent.receiverAddress = controller.provider.ConvertPubKeyToAddress(receiverPubkey)
		typedEvent.senderAddress = event.Address
		typedEvents = append(typedEvents, typedEvent)
	}

	// Then, handle multi transfers
	for _, event := range rawEventsMultiDCDTNFTTransfer {
		numTopics := len(event.Topics)
		numTopicsExceptLast := numTopics - 1
		numTopicsPerTransfer := numTopicsPerTransferOfEventMultiDCDTNFTTransfer

		if numTopicsExceptLast%numTopicsPerTransfer != 0 {
			return nil, fmt.Errorf("%w: bad number of topics for MultiDCDTNFTTransfer event = %d", errCannotRecognizeEvent, numTopics)
		}

		numTransfers := numTopicsExceptLast / numTopicsPerTransfer
		receiverPubkey := event.Topics[numTopics-1]
		receiver := controller.provider.ConvertPubKeyToAddress(receiverPubkey)

		for i := 0; i < numTransfers; i++ {
			typedEvent, err := newEventDCDTFromBasicTopics(event.Topics[i*numTopicsPerTransfer+0 : i*numTopicsPerTransfer+3])
			if err != nil {
				return nil, err
			}

			typedEvent.receiverAddress = receiver
			typedEvent.senderAddress = event.Address
			typedEvents = append(typedEvents, typedEvent)
		}
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsDCDTLocalBurn(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventDCDTLocalBurn)
	typedEvents := make([]*eventDCDT, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTLocalBurn {
			return nil, fmt.Errorf("%w: bad number of topics for DCDTLocalBurn event = %d", errCannotRecognizeEvent, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		typedEvent.otherAddress = event.Address
		typedEvents = append(typedEvents, typedEvent)
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsDCDTLocalMint(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventDCDTLocalMint)
	typedEvents := make([]*eventDCDT, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTLocalMint {
			return nil, fmt.Errorf("%w: bad number of topics for DCDTLocalMint event = %d", errCannotRecognizeEvent, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		typedEvent.otherAddress = event.Address
		typedEvents = append(typedEvents, typedEvent)
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsDCDTWipe(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventDCDTWipe)
	typedEvents := make([]*eventDCDT, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTWipe {
			return nil, fmt.Errorf("%w: bad number of topics for DCDTWipe event = %d", errCannotRecognizeEvent, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		accountPubkey := event.Topics[3]
		typedEvent.otherAddress = controller.provider.ConvertPubKeyToAddress(accountPubkey)
		typedEvents = append(typedEvents, typedEvent)
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsDCDTNFTCreate(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventDCDTNFTCreate)
	typedEvents := make([]*eventDCDT, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTNFTCreate {
			return nil, fmt.Errorf("%w: bad number of topics for %s event = %d", errCannotRecognizeEvent, transactionEventDCDTNFTCreate, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		typedEvent.otherAddress = event.Address
		typedEvents = append(typedEvents, typedEvent)
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsDCDTNFTBurn(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventDCDTNFTBurn)
	typedEvents := make([]*eventDCDT, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTNFTBurn {
			return nil, fmt.Errorf("%w: bad number of topics for %s event = %d", errCannotRecognizeEvent, transactionEventDCDTNFTBurn, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		typedEvent.otherAddress = event.Address
		typedEvents = append(typedEvents, typedEvent)
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsDCDTNFTAddQuantity(tx *transaction.ApiTransactionResult) ([]*eventDCDT, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventDCDTNFTAddQuantity)
	typedEvents := make([]*eventDCDT, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventDCDTNFTAddQuantity {
			return nil, fmt.Errorf("%w: bad number of topics for %s event = %d", errCannotRecognizeEvent, transactionEventDCDTNFTAddQuantity, numTopics)
		}

		typedEvent, err := newEventDCDTFromBasicTopics(event.Topics)
		if err != nil {
			return nil, err
		}

		typedEvent.otherAddress = event.Address
		typedEvents = append(typedEvents, typedEvent)
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) extractEventsClaimDeveloperRewards(tx *transaction.ApiTransactionResult) ([]*eventClaimDeveloperRewards, error) {
	rawEvents := controller.findManyEventsByIdentifier(tx, transactionEventClaimDeveloperRewards)
	typedEvents := make([]*eventClaimDeveloperRewards, 0, len(rawEvents))

	for _, event := range rawEvents {
		numTopics := len(event.Topics)
		if numTopics != numTopicsOfEventClaimDeveloperRewards {
			return nil, fmt.Errorf("%w: bad number of topics for %s event = %d", errCannotRecognizeEvent, transactionEventClaimDeveloperRewards, numTopics)
		}

		valueBytes := event.Topics[0]
		receiverPubkey := event.Topics[1]

		value := big.NewInt(0).SetBytes(valueBytes)
		receiver := controller.provider.ConvertPubKeyToAddress(receiverPubkey)

		typedEvents = append(typedEvents, &eventClaimDeveloperRewards{
			value:           value.String(),
			receiverAddress: receiver,
		})
	}

	return typedEvents, nil
}

func (controller *transactionEventsController) findManyEventsByIdentifier(tx *transaction.ApiTransactionResult, identifier string) []*transaction.Events {
	events := make([]*transaction.Events, 0)

	if !controller.hasEvents(tx) {
		return events
	}

	for _, event := range tx.Logs.Events {
		if event.Identifier == identifier {
			events = append(events, event)
		}
	}

	return events
}

func (controller *transactionEventsController) hasSignalErrorOfMetaTransactionIsInvalid(tx *transaction.ApiTransactionResult) bool {
	if !controller.hasEvents(tx) {
		return false
	}

	for _, event := range tx.Logs.Events {
		isSignalError := event.Identifier == transactionEventSignalError
		if !isSignalError {
			continue
		}

		if eventHasTopic(event, transactionEventTopicInvalidMetaTransaction) {
			return true
		}
		if eventHasTopic(event, transactionEventTopicInvalidMetaTransactionNotEnoughGas) {
			return true
		}
	}

	return false
}

func (controller *transactionEventsController) hasEvents(tx *transaction.ApiTransactionResult) bool {
	return tx.Logs != nil && tx.Logs.Events != nil && len(tx.Logs.Events) > 0
}

func eventHasTopic(event *transaction.Events, topicToFind string) bool {
	for _, topic := range event.Topics {
		if string(topic) == topicToFind {
			return true
		}
	}

	return false
}
