/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resmgmt

import (
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/pkg/errors"
)

// WithTargets allows overriding of the target peers for the request.
func WithTargets(targets ...fab.Peer) RequestOption {
	return func(ctx context.Client, opts *requestOptions) error {
		opts.Targets = targets
		return nil
	}
}

// WithTargetURLs allows overriding of the target peers for the request.
// Targets are specified by URL, and the SDK will create the underlying peer
// objects.
func WithTargetURLs(urls ...string) RequestOption {
	return func(ctx context.Client, opts *requestOptions) error {

		var targets []fab.Peer

		for _, url := range urls {

			peerCfg, err := config.NetworkPeerConfigFromURL(ctx.Config(), url)
			if err != nil {
				return err
			}

			peer, err := ctx.InfraProvider().CreatePeerFromConfig(peerCfg)
			if err != nil {
				return errors.WithMessage(err, "creating peer from config failed")
			}

			targets = append(targets, peer)
		}

		return WithTargets(targets...)(ctx, opts)
	}
}

// WithTargetFilter enables a target filter for the request.
func WithTargetFilter(targetFilter TargetFilter) RequestOption {
	return func(ctx context.Client, opts *requestOptions) error {
		opts.TargetFilter = targetFilter
		return nil
	}
}

// WithTimeout specifies the timeout for the request. If not specified,
// a default timeout from configuration is used.
func WithTimeout(timeout time.Duration) RequestOption {
	return func(ctx context.Client, opts *requestOptions) error {
		opts.Timeout = timeout
		return nil
	}
}

// WithOrdererURL allows an orderer to be specified for the request.
// The orderer will be looked-up based on the url argument.
// A default orderer implementation will be used.
func WithOrdererURL(url string) RequestOption {
	return func(ctx context.Client, opts *requestOptions) error {

		ordererCfg, err := ctx.Config().OrdererConfig(url)
		if err != nil {
			return errors.WithMessage(err, "orderer not found")
		}
		if ordererCfg == nil {
			return errors.New("orderer not found")
		}

		orderer, err := ctx.InfraProvider().CreateOrdererFromConfig(ordererCfg)
		if err != nil {
			return errors.WithMessage(err, "creating orderer from config failed")
		}

		return WithOrderer(orderer)(ctx, opts)
	}
}

// WithOrderer allows an orderer to be specified for the request.
func WithOrderer(orderer fab.Orderer) RequestOption {
	return func(ctx context.Client, opts *requestOptions) error {
		opts.Orderer = orderer
		return nil
	}
}
