//go:build boringcrypto

package main

import _ "crypto/tls/fipsonly"

// The fipsonly import enforces strict FIPS 140-3 compliance
// This requires Extended Master Secret (EMS) for TLS 1.2 connections
// If a server doesn't support EMS, the connection will fail with:
// "tls: FIPS 140-3 requires the use of Extended Master Secret"
// This is the expected and correct behavior for FIPS 140-3 compliance

func checkBoringCrypto() bool {
	return true
}
