/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabsdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/msp"
	"github.com/pkg/errors"
)

type identityOptions struct {
	identity msp.Identity
	orgName  string
	userName string
}

// ContextOption provides parameters for creating a session (primarily from a fabric identity/user)
type ContextOption func(s *identityOptions) error

// WithUser uses the named user to load the identity
func WithUser(userName string) ContextOption {
	return func(o *identityOptions) error {
		o.userName = userName
		return nil
	}
}

// WithIdentity uses a pre-constructed identity object as the credential for the session
func WithIdentity(identity msp.Identity) ContextOption {
	return func(o *identityOptions) error {
		o.identity = identity
		return nil
	}
}

// WithOrg uses the named organization
func WithOrg(org string) ContextOption {
	return func(o *identityOptions) error {
		o.orgName = org
		return nil
	}
}

// ErrAnonymousIdentity is returned when options for identity creation
// don't include neither username nor identity
var ErrAnonymousIdentity = errors.New("missing credentials")

func (sdk *FabricSDK) newIdentity(options ...ContextOption) (msp.Identity, error) {
	clientConfig, err := sdk.Config().Client()
	if err != nil {
		return nil, errors.WithMessage(err, "retrieving client configuration failed")
	}

	opts := identityOptions{
		orgName: clientConfig.Organization,
	}

	for _, option := range options {
		err := option(&opts)
		if err != nil {
			return nil, errors.WithMessage(err, "error in option passed to create identity")
		}
	}

	if opts.identity == nil && opts.userName == "" {
		return nil, ErrAnonymousIdentity
	}

	if opts.identity != nil {
		return opts.identity, nil
	}

	if opts.userName == "" || opts.orgName == "" {
		return nil, errors.New("invalid options to create identity")
	}

	mgr, ok := sdk.provider.IdentityManager(opts.orgName)
	if !ok {
		return nil, errors.New("invalid options to create identity, invalid org name")
	}

	user, err := mgr.GetUser(opts.userName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// session represents an identity being used with clients along with services
// that associate with that identity (particularly the channel service).
type session struct {
	msp.Identity
}

// newSession creates a session from a context and a user (TODO)
func newSession(ic msp.Identity, cp fab.ChannelProvider) *session {
	s := session{
		Identity: ic,
	}

	return &s
}
