// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vanadium/services/identity/internal/oauth"
	"github.com/vanadium/services/identity/internal/revocation"
	"github.com/vanadium/services/identity/internal/util"
	v23 "v.io/v23"
	"v.io/v23/context"
	"v.io/v23/security"
	"v.io/v23/vom"
	"v.io/x/ref/lib/stats"
	"v.io/x/ref/lib/stats/counter"
)

const (
	publicKeyFormKey    = "public_key"
	tokenFormKey        = "token"
	caveatsFormKey      = "caveats"
	outputFormatFormKey = "output_format"

	jsonFormat      = "json"
	base64VomFormat = "base64vom"
)

// Map of client id -> blessing extension
// The blessing given for a token with ClientID 'id' is generated by extending
// the default blessing of this service's principal with the extension held
// in this map for 'id'.
// The string "{email}" in the Extension will be replaced by the email from the
// request's access token.
type RegisteredAppMap map[string]struct {
	Extension string
}

// OAuthBlesserParams represents all the parameters required for exchanging an
// OAuth token for blessings.
type OAuthBlesserParams struct {
	// The OAuth provider that must have issued the access tokens accepted by ths service.
	OAuthProvider oauth.OAuthProvider
	// The object name of the discharger service. If this is empty then revocation caveats will not be granted.
	DischargerLocation string
	// The revocation manager that generates caveats and manages revocation.
	RevocationManager revocation.RevocationManager
	// The duration for which blessings will be valid. (Used iff RevocationManager is nil).
	BlessingDuration time.Duration
}

type accessTokenBlesser struct {
	ctx    *context.T
	params OAuthBlesserParams
	apps   RegisteredAppMap

	countersMu sync.Mutex
	counters   map[string]*counter.Counter
}

// NewOAuthBlessingHandler returns an http.Handler that uses Google OAuth2 Access tokens
// to obtain the username of the requestor and reponds with blessings for that username.
//
// The blessings are namespaced under the ClientID for the access token. In particular,
// the name of the granted blessing is of the form <idp>:<appID>:<email> where <idp>
// is the name of the default blessings used by the identity provider and
// <appID> is the name of the 'app' - either the OAuth ClientID or a registered
// alias.
//
// Blessings generated by this service carry a third-party revocation caveat if a
// RevocationManager is specified by the params or they carry an ExpiryCaveat that
// expires after the duration specified by the params.
//
// The handler expects the following request parameters:
//   - "public_key": Base64 DER encoded PKIX representation of the client's public key
//   - "caveats": Base64 VOM encoded list of caveats [OPTIONAL]
//   - "token": Google OAuth2 Access token
//   - "output_format": The encoding format for the returned blessings. The following
//     formats are supported:
//   - "json": JSON-encoding of the wire format of Blessings.
//   - "base64vom": Base64URL encoding of VOM-encoded Blessings [DEFAULT]
//
// The response consists of blessings encoded in the requested output format.
//
// WARNINGS:
//   - There is no binding between the channel over which the access token
//     was obtained and the channel used to make this request.
//   - There is no "proof of possession of private key" required by the server.
//
// Thus, if Mallory (attacker) possesses the access token associated with Alice's
// account (victim), she may be able to obtain a blessing with Alice's name on it
// for any public key of her choice.
func NewOAuthBlessingHandler(ctx *context.T, params OAuthBlesserParams, apps RegisteredAppMap) http.Handler {
	if apps == nil {
		apps = make(RegisteredAppMap)
	}
	return &accessTokenBlesser{ctx: ctx, params: params, apps: apps}
}

func (a *accessTokenBlesser) blessingCaveats(r *http.Request, p security.Principal) ([]security.Caveat, error) {
	var caveats []security.Caveat
	if base64VomCaveats := r.FormValue(caveatsFormKey); len(base64VomCaveats) != 0 {
		vomCaveats, err := base64.URLEncoding.DecodeString(base64VomCaveats)
		if err != nil {
			return nil, fmt.Errorf("base64.URLEncoding.DecodeString failed: %v", err)
		}

		if err := vom.Decode(vomCaveats, &caveats); err != nil {
			return nil, fmt.Errorf("vom.Decode failed: %v", err)
		}
	}
	if len(caveats) == 0 {
		var (
			cav security.Caveat
			err error
		)
		if a.params.RevocationManager != nil {
			cav, err = a.params.RevocationManager.NewCaveat(p.PublicKey(), a.params.DischargerLocation)
		} else {
			cav, err = security.NewExpiryCaveat(time.Now().Add(a.params.BlessingDuration))
		}
		if err != nil {
			return nil, fmt.Errorf("failed to construct caveats: %v", err)
		}
		caveats = append(caveats, cav)
	}
	return caveats, nil

}

func (a *accessTokenBlesser) remotePublicKey(r *http.Request) (security.PublicKey, error) {
	publicKeyVom, err := base64.URLEncoding.DecodeString(r.FormValue(publicKeyFormKey))
	if err != nil {
		return nil, fmt.Errorf("base64.URLEncoding.DecodeString failed: %v", err)
	}
	return security.UnmarshalPublicKey(publicKeyVom)
}

func (a *accessTokenBlesser) blessingExtension(r *http.Request) (string, error) {
	email, clientID, err := a.params.OAuthProvider.GetEmailAndClientID(r.FormValue(tokenFormKey))
	if err != nil {
		return "", err
	}
	if entry, ok := a.apps[clientID]; ok {
		return strings.ReplaceAll(entry.Extension, "{email}", email), nil
	}

	// We use <clientID>/<email> as the extension in order to namespace the blessing under
	// the <clientID>. This has the downside that the blessing cannot be used to act on
	// behalf of the user, i.e., services access controlled to blessings matching"<idp>/<email>"
	// would not authorize this blessing.
	//
	// The alternative is to use the extension <email>/<clientID> however this is risky as it
	// may provide too much authority to the app, especially since we don't have a set of default
	// caveats to apply to the blessing.
	//
	// TODO(ataly, ashankar): Think about changing to the extension <email>/<clientID>.
	return strings.Join([]string{clientID, email}, security.ChainSeparator), nil
}

func (a *accessTokenBlesser) encodeBlessingsJSON(b security.Blessings) ([]byte, error) {
	return json.Marshal(security.MarshalBlessings(b))
}

func (a *accessTokenBlesser) encodeBlessingsVOM(b security.Blessings) (string, error) {
	bVom, err := vom.Encode(b)
	if err != nil {
		return "", fmt.Errorf("vom.Encode(%v) failed: %v", b, err)
	}
	return base64.URLEncoding.EncodeToString(bVom), nil
}

func (a *accessTokenBlesser) counter(name string) *counter.Counter {
	a.countersMu.Lock()
	defer a.countersMu.Unlock()
	if a.counters == nil {
		a.counters = make(map[string]*counter.Counter)
	}
	if c := a.counters[name]; c != nil {
		return c
	}
	c := stats.NewCounter(fmt.Sprint("http/blessings/", name))
	a.counters[name] = c
	return c
}

func (a *accessTokenBlesser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteKey, err := a.remotePublicKey(r)
	if err != nil {
		a.ctx.Infof("Failed to decode public key [%v] for request %#v", err, r)
		util.HTTPServerError(w, fmt.Errorf("failed to decode public key: %v", err))
		return
	}

	p := v23.GetPrincipal(a.ctx)
	with, _ := p.BlessingStore().Default()

	caveats, err := a.blessingCaveats(r, p)
	if err != nil {
		a.ctx.Infof("Failed to constuct caveats for blessing [%v] for request %#v", err, r)
		util.HTTPServerError(w, fmt.Errorf("failed to construct caveats for blessing: %v", err))
		return
	}

	extension, err := a.blessingExtension(r)
	if err != nil {
		a.ctx.Infof("Failed to process access token [%v] for request %#v", err, r)
		util.HTTPServerError(w, fmt.Errorf("failed to process access token: %v", err))
		return
	}

	blessings, err := p.Bless(remoteKey, with, extension, caveats[0], caveats[1:]...)
	if err != nil {
		a.ctx.Infof("Failed to Bless [%v] for request %#v", err, r)
		util.HTTPServerError(w, fmt.Errorf("failed to Bless: %v", err))
		return
	}
	a.counter(with.String()).Incr(1)

	outputFormat := r.FormValue(outputFormatFormKey)
	if len(outputFormat) == 0 {
		outputFormat = base64VomFormat
	}
	switch outputFormat {
	case jsonFormat:
		encodedBlessings, err := a.encodeBlessingsJSON(blessings)
		if err != nil {
			a.ctx.Infof("Failed to encode blessings [%v] for request %#v", err, r)
			util.HTTPServerError(w, fmt.Errorf("failed to encode blessings in format %v: %v", outputFormat, err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedBlessings) //nolint:errcheck
	case base64VomFormat:
		encodedBlessings, err := a.encodeBlessingsVOM(blessings)
		if err != nil {
			a.ctx.Infof("Failed to encode blessings [%v] for request %#v", err, r)
			util.HTTPServerError(w, fmt.Errorf("failed to encode blessings in format %v: %v", outputFormat, err))
			return
		}
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte(encodedBlessings)) //nolint:errcheck
	default:
		a.ctx.Infof("Unrecognized output format [%v] in request %#v", outputFormat, r)
		util.HTTPServerError(w, fmt.Errorf("unrecognized output format [%v] in request. Allowed formats are [%v, %v]", outputFormat, base64VomFormat, jsonFormat))
		return
	}
}
