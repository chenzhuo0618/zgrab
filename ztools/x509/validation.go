// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

import "time"

// Validation stores different validation levels for a given certificate
type Validation struct {
	BrowserTrusted bool   `json:"browser_trusted"`
	BrowserError   string `json:"browser_error,omitempty"`
	MatchesDomain  bool   `json:"-"`
	Domain         string `json:"-"`
}

// ValidateWithStupidDetail fills out a Validation struct given a leaf
// certificate and intermediates / roots. If opts.DNSName is set, then it will
// also check if the domain matches.
func (c *Certificate) ValidateWithStupidDetail(opts VerifyOptions) *Validation {

	// Manually set the time, so that all verifies we do get the same time
	if opts.CurrentTime.IsZero() {
		opts.CurrentTime = time.Now()
	}

	opts.KeyUsages = nil

	out := new(Validation)
	out.Domain = opts.DNSName

	if _, err := c.Verify(opts); err != nil {
		switch err := err.(type) {
		case HostnameError:
			out.BrowserTrusted = true
			out.MatchesDomain = false
		default:
			out.BrowserTrusted = false
			out.BrowserError = err.Error()
		}
	} else {
		out.BrowserTrusted = true
		if opts.DNSName != "" {
			out.MatchesDomain = true
		}
	}

	return out
}
