/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"fmt"

	reqContext "context"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
)

// MockInfraProvider represents the default implementation of Fabric objects.
type MockInfraProvider struct {
	providerContext  context.Providers
	customOrderer    fab.Orderer
	customTransactor fab.Transactor
}

// CreateEventService creates the event service.
func (f *MockInfraProvider) CreateEventService(ic fab.ClientContext, chConfig fab.ChannelCfg) (fab.EventService, error) {
	panic("not implemented")
}

// CreateChannelConfig initializes the channel config
func (f *MockInfraProvider) CreateChannelConfig(channelID string) (fab.ChannelConfig, error) {
	return &MockChannelConfig{channelID: channelID}, nil
}

// CreateChannelMembership returns a channel member identifier
func (f *MockInfraProvider) CreateChannelMembership(cfg fab.ChannelCfg) (fab.ChannelMembership, error) {
	return nil, fmt.Errorf("Not implemented")
}

// CreateChannelTransactor initializes the transactor
func (f *MockInfraProvider) CreateChannelTransactor(reqCtx reqContext.Context, cfg fab.ChannelCfg) (fab.Transactor, error) {
	if f.customTransactor != nil {
		return f.customTransactor, nil
	}
	if cfg == nil {
		return &MockTransactor{}, nil
	}
	return &MockTransactor{ChannelID: cfg.ID(), Ctx: reqCtx}, nil
}

// CreatePeerFromConfig returns a new default implementation of Peer based configuration
func (f *MockInfraProvider) CreatePeerFromConfig(peerCfg *core.NetworkPeer) (fab.Peer, error) {
	if peerCfg != nil {
		p := NewMockPeer(peerCfg.MspID, peerCfg.URL)
		p.SetMSPID(peerCfg.MspID)

		return p, nil
	}
	return &MockPeer{}, nil
}

// CreateOrdererFromConfig creates a default implementation of Orderer based on configuration.
func (f *MockInfraProvider) CreateOrdererFromConfig(cfg *core.OrdererConfig) (fab.Orderer, error) {
	if f.customOrderer != nil {
		return f.customOrderer, nil
	}

	return &MockOrderer{}, nil
}

//CommManager returns comm provider
func (f *MockInfraProvider) CommManager() fab.CommManager {
	return nil
}

// SetCustomOrderer creates a default implementation of Orderer based on configuration.
func (f *MockInfraProvider) SetCustomOrderer(customOrderer fab.Orderer) {
	f.customOrderer = customOrderer
}

// SetCustomTransactor sets custom transactor for unit-test purposes
func (f *MockInfraProvider) SetCustomTransactor(customTransactor fab.Transactor) {
	f.customTransactor = customTransactor
}

//Close mock close function
func (f *MockInfraProvider) Close() {
}
