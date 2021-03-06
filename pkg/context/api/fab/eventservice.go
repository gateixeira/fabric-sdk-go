/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fab

import (
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

// BlockEvent contains the data for the block event
type BlockEvent struct {
	Block *cb.Block
}

// FilteredBlockEvent contains the data for a filtered block event
type FilteredBlockEvent struct {
	FilteredBlock *pb.FilteredBlock
}

// TxStatusEvent contains the data for a transaction status event
type TxStatusEvent struct {
	TxID             string
	TxValidationCode pb.TxValidationCode
}

// CCEvent contains the data for a chaincode event
type CCEvent struct {
	TxID        string
	ChaincodeID string
	EventName   string
}

// Registration is a handle that is returned from a successful RegisterXXXEvent.
// This handle should be used in Unregister in order to unregister the event.
type Registration interface{}

// BlockFilter is a function that determines whether a Block event
// should be ignored
type BlockFilter func(block *cb.Block) bool

// EventService is a service that receives events such as block, filtered block,
// chaincode, and transaction status events.
type EventService interface {
	// RegisterBlockEvent registers for block events. If the caller does not have permission
	// to register for block events then an error is returned.
	// Note that Unregister must be called when the registration is no longer needed.
	// - filter is an optional filter that filters out unwanted events. (Note: Only one filter may be specified.)
	// - Returns the registration and a channel that is used to receive events. The channel
	//   is closed when Unregister is called.
	RegisterBlockEvent(filter ...BlockFilter) (Registration, <-chan *BlockEvent, error)

	// RegisterFilteredBlockEvent registers for filtered block events.
	// Note that Unregister must be called when the registration is no longer needed.
	// - Returns the registration and a channel that is used to receive events. The channel
	//   is closed when Unregister is called.
	RegisterFilteredBlockEvent() (Registration, <-chan *FilteredBlockEvent, error)

	// RegisterChaincodeEvent registers for chaincode events.
	// Note that Unregister must be called when the registration is no longer needed.
	// - ccID is the chaincode ID for which events are to be received
	// - eventFilter is the chaincode event filter (regular expression) for which events are to be received
	// - Returns the registration and a channel that is used to receive events. The channel
	//   is closed when Unregister is called.
	RegisterChaincodeEvent(ccID, eventFilter string) (Registration, <-chan *CCEvent, error)

	// RegisterTxStatusEvent registers for transaction status events.
	// Note that Unregister must be called when the registration is no longer needed.
	// - txID is the transaction ID for which events are to be received
	// - Returns the registration and a channel that is used to receive events. The channel
	//   is closed when Unregister is called.
	RegisterTxStatusEvent(txID string) (Registration, <-chan *TxStatusEvent, error)

	// Unregister removes the given registration and closes the event channel.
	// - reg is the registration handle that was returned from one of the Register functions
	Unregister(reg Registration)
}

// ConnectionEvent is sent when the client disconnects from or
// reconnects to the event server. Connected == true means that the
// client has connected, whereas Connected == false means that the
// client has disconnected. In the disconnected case, Err contains
// the disconnect error.
type ConnectionEvent struct {
	Connected bool
	Err       error
}

// EventClient is a client that connects to a peer and receives channel events
// such as block, filtered block, chaincode, and transaction status events.
type EventClient interface {
	EventService

	// Connect connects to the event server.
	Connect() error

	// Close closes the connection to the event server and releases all resources.
	// Once this function is invoked the client may no longer be used.
	Close()

	// CloseIfIdle closes the connection to the event server only if there are no outstanding
	// registrations.
	// Returns true if the client was closed. In this case the client may no longer be used.
	// A return value of false indicates that the client could not be closed since
	// there was at least one registration.
	CloseIfIdle() bool
}
