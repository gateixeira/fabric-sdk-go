/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defcore"
)

// ========== MSP Provider Factory with custom crypto provider ============= //

// CustomCoreFactory is a custom factory for tests.
type CustomCoreFactory struct {
	defaultFactory    *defcore.ProviderFactory
	customCryptoSuite core.CryptoSuite
}

// NewCustomCoreFactory creates a new instance of customCoreFactory
func NewCustomCoreFactory(customCryptoSuite core.CryptoSuite) *CustomCoreFactory {
	return &CustomCoreFactory{defaultFactory: defcore.NewProviderFactory(), customCryptoSuite: customCryptoSuite}
}

// CreateCryptoSuiteProvider creates custom crypto provider
func (f *CustomCoreFactory) CreateCryptoSuiteProvider(config core.Config) (core.CryptoSuite, error) {
	return f.customCryptoSuite, nil
}

// CreateSigningManager creates SigningManager
func (f *CustomCoreFactory) CreateSigningManager(cryptoProvider core.CryptoSuite, config core.Config) (core.SigningManager, error) {
	return f.defaultFactory.CreateSigningManager(cryptoProvider, config)
}

// CreateInfraProvider creates InfraProvider
func (f *CustomCoreFactory) CreateInfraProvider(config core.Config) (fab.InfraProvider, error) {
	return f.defaultFactory.CreateInfraProvider(config)
}
