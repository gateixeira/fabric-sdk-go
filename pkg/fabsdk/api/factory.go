/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/msp"
)

// Providers represents the SDK configured providers context.
type Providers interface {
	core.Providers
	msp.Providers
	fab.Providers
}

// CoreProviderFactory allows overriding of primitives and the fabric core object provider
type CoreProviderFactory interface {
	CreateCryptoSuiteProvider(config core.Config) (core.CryptoSuite, error)
	CreateSigningManager(cryptoProvider core.CryptoSuite, config core.Config) (core.SigningManager, error)
	CreateInfraProvider(config core.Config) (fab.InfraProvider, error)
}

// MSPProviderFactory allows overriding providers of MSP services
type MSPProviderFactory interface {
	CreateUserStore(config core.Config) (msp.UserStore, error)
	CreateProvider(config core.Config, cryptoProvider core.CryptoSuite, userStore msp.UserStore) (msp.Provider, error)
}

// ServiceProviderFactory allows overriding default service providers (such as peer discovery)
type ServiceProviderFactory interface {
	CreateDiscoveryProvider(config core.Config, fabPvdr fab.InfraProvider) (fab.DiscoveryProvider, error)
	CreateSelectionProvider(config core.Config) (fab.SelectionProvider, error)
}
