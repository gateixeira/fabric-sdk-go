/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context

import (
	reqContext "context"

	"github.com/pkg/errors"

	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/msp"
)

// Client supplies the configuration and signing identity to client objects.
type Client struct {
	context.Providers
	msp.Identity
}

//Channel supplies the configuration for channel context client
type Channel struct {
	context.Client
	discovery      fab.DiscoveryService
	selection      fab.SelectionService
	channelService fab.ChannelService
	channelID      string
}

//Providers returns core providers
func (c *Channel) Providers() context.Client {
	return c
}

//DiscoveryService returns core discovery service
func (c *Channel) DiscoveryService() fab.DiscoveryService {
	return c.discovery
}

//SelectionService returns selection service
func (c *Channel) SelectionService() fab.SelectionService {
	return c.selection
}

//ChannelService returns channel service
func (c *Channel) ChannelService() fab.ChannelService {
	return c.channelService
}

//ChannelID returns channel id
func (c *Channel) ChannelID() string {
	return c.channelID
}

//Provider implementation for Providers interface
type Provider struct {
	config            core.Config
	userStore         msp.UserStore
	cryptoSuite       core.CryptoSuite
	discoveryProvider fab.DiscoveryProvider
	selectionProvider fab.SelectionProvider
	signingManager    core.SigningManager
	mspProvider       msp.Provider
	infraProvider     fab.InfraProvider
	channelProvider   fab.ChannelProvider
}

// Config returns the Config provider of sdk.
func (c *Provider) Config() core.Config {
	return c.config
}

// CryptoSuite returns the BCCSP provider of sdk.
func (c *Provider) CryptoSuite() core.CryptoSuite {
	return c.cryptoSuite
}

// IdentityManager returns identity manager for organization
func (c *Provider) IdentityManager(orgName string) (msp.IdentityManager, bool) {
	return c.mspProvider.IdentityManager(orgName)
}

// SigningManager returns signing manager
func (c *Provider) SigningManager() core.SigningManager {
	return c.signingManager
}

// UserStore returns state store
func (c *Provider) UserStore() msp.UserStore {
	return c.userStore
}

// DiscoveryProvider returns discovery provider
func (c *Provider) DiscoveryProvider() fab.DiscoveryProvider {
	return c.discoveryProvider
}

// SelectionProvider returns selection provider
func (c *Provider) SelectionProvider() fab.SelectionProvider {
	return c.selectionProvider
}

// ChannelProvider provides channel services.
func (c *Provider) ChannelProvider() fab.ChannelProvider {
	return c.channelProvider
}

// InfraProvider provides fabric objects such as peer and user
func (c *Provider) InfraProvider() fab.InfraProvider {
	return c.infraProvider
}

//SDKContextParams parameter for creating FabContext
type SDKContextParams func(opts *Provider)

//WithConfig sets config to FabContext
func WithConfig(config core.Config) SDKContextParams {
	return func(ctx *Provider) {
		ctx.config = config
	}
}

// WithUserStore sets user store to FabContext
func WithUserStore(userStore msp.UserStore) SDKContextParams {
	return func(ctx *Provider) {
		ctx.userStore = userStore
	}
}

//WithCryptoSuite sets cryptosuite parameter to FabContext
func WithCryptoSuite(cryptoSuite core.CryptoSuite) SDKContextParams {
	return func(ctx *Provider) {
		ctx.cryptoSuite = cryptoSuite
	}
}

//WithDiscoveryProvider sets discoveryProvider to FabContext
func WithDiscoveryProvider(discoveryProvider fab.DiscoveryProvider) SDKContextParams {
	return func(ctx *Provider) {
		ctx.discoveryProvider = discoveryProvider
	}
}

//WithSelectionProvider sets selectionProvider to FabContext
func WithSelectionProvider(selectionProvider fab.SelectionProvider) SDKContextParams {
	return func(ctx *Provider) {
		ctx.selectionProvider = selectionProvider
	}
}

//WithSigningManager sets signingManager to FabContext
func WithSigningManager(signingManager core.SigningManager) SDKContextParams {
	return func(ctx *Provider) {
		ctx.signingManager = signingManager
	}
}

//WithMSPProvider sets MSPProvider maps to context
func WithMSPProvider(provider msp.Provider) SDKContextParams {
	return func(ctx *Provider) {
		ctx.mspProvider = provider
	}
}

//WithInfraProvider sets infraProvider maps to FabContext
func WithInfraProvider(infraProvider fab.InfraProvider) SDKContextParams {
	return func(ctx *Provider) {
		ctx.infraProvider = infraProvider
	}
}

//WithChannelProvider sets channelProvider to FabContext
func WithChannelProvider(channelProvider fab.ChannelProvider) SDKContextParams {
	return func(ctx *Provider) {
		ctx.channelProvider = channelProvider
	}
}

//NewProvider creates new context client provider
// Not be used by end developers, fabsdk package use only
func NewProvider(params ...SDKContextParams) *Provider {
	ctxProvider := Provider{}
	for _, param := range params {
		param(&ctxProvider)
	}
	return &ctxProvider
}

//NewChannel creates new channel context client
// Not be used by end developers, fabsdk package use only
func NewChannel(clientProvider context.ClientProvider, channelID string) (*Channel, error) {

	client, err := clientProvider()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get client context to create channel client")
	}

	channelService, err := client.ChannelProvider().ChannelService(client, channelID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get channel service to create channel client")
	}

	discoveryService, err := client.DiscoveryProvider().CreateDiscoveryService(channelID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get discovery service to create channel client")
	}

	selectionService, err := client.SelectionProvider().CreateSelectionService(channelID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get selection service to create channel client")
	}

	return &Channel{
		Client:         client,
		selection:      selectionService,
		discovery:      discoveryService,
		channelService: channelService,
		channelID:      channelID,
	}, nil
}

type reqContextKey string

var reqContextCommManager = reqContextKey("commManager")
var reqContextTimeout = reqContextKey("timeout")
var reqContextClient = reqContextKey("clientContext")

//WithTimeoutType sets timeout by type defined in config to request context
func WithTimeoutType(timeoutType core.TimeoutType) ReqContextOptions {
	return func(ctx *requestContextOpts) {
		ctx.timeoutType = timeoutType
	}
}

//WithTimeout sets timeout time duration to request context
func WithTimeout(timeout time.Duration) ReqContextOptions {
	return func(ctx *requestContextOpts) {
		ctx.timeout = timeout
	}
}

//WithReqContext sets existing reqContext as a parent ReqContext
func WithReqContext(context reqContext.Context) ReqContextOptions {
	return func(ctx *requestContextOpts) {
		ctx.parentContext = context
	}
}

//ReqContextOptions parameter for creating requestContext
type ReqContextOptions func(opts *requestContextOpts)

type requestContextOpts struct {
	timeoutType   core.TimeoutType
	timeout       time.Duration
	parentContext reqContext.Context
}

// NewRequest creates a request-scoped context.
func NewRequest(client context.Client, options ...ReqContextOptions) (reqContext.Context, reqContext.CancelFunc) {

	//'-1' to get default config timeout when timeout options not passed
	reqCtxOpts := requestContextOpts{timeoutType: -1}
	for _, option := range options {
		option(&reqCtxOpts)
	}

	var timeout time.Duration
	if reqCtxOpts.timeout > 0 {
		timeout = reqCtxOpts.timeout
	} else {
		timeout = client.Config().TimeoutOrDefault(reqCtxOpts.timeoutType)
	}

	parentContext := reqCtxOpts.parentContext
	if parentContext == nil {
		//when parent request context not set, use background context
		parentContext = reqContext.Background()
	}
	ctx := reqContext.WithValue(parentContext, reqContextCommManager, client.InfraProvider().CommManager())
	ctx = reqContext.WithValue(ctx, reqContextClient, client)
	ctx, cancel := reqContext.WithTimeout(ctx, timeout)

	return ctx, cancel
}

// RequestCommManager extracts the CommManager from the request-scoped context.
func RequestCommManager(ctx reqContext.Context) (fab.CommManager, bool) {
	commManager, ok := ctx.Value(reqContextCommManager).(fab.CommManager)
	return commManager, ok
}

// RequestClientContext extracts the Client Context from the request-scoped context.
func RequestClientContext(ctx reqContext.Context) (context.Client, bool) {
	clientContext, ok := ctx.Value(reqContextClient).(context.Client)
	return clientContext, ok
}
