// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package restsigner

import (
	"fmt"

	"v.io/v23/security"
)

func NewRestSigner() (security.Signer, error) {
	return nil, fmt.Errorf("restsigner is currently not implemented")
}
