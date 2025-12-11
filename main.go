package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Note: In Go 1.24+, FIPS mode uses two environment variables:
// - GOFIPS140=v1.0.0 (or =latest): Build-time - selects FIPS module version
// - GODEBUG=fips140=on: Runtime - enables FIPS mode execution
// - GODEBUG=fips140=only: Runtime - strict FIPS enforcement (panics on non-FIPS crypto)

func main() {
	// Check runtime FIPS mode control
	godebug := os.Getenv("GODEBUG")
	fipsMode := getFIPSMode(godebug)

	fmt.Println("=== FIPS Mode Status ===")
	fmt.Println("")
	fmt.Println("ℹ️  FIPS enforcement is determined at BUILD time with GOFIPS140=v1.0.0")
	fmt.Println("ℹ️  When built with GOFIPS140, FIPS 140-3 compliance is ALWAYS enforced")
	fmt.Println("")

	if fipsMode == "only" {
		fmt.Println("GODEBUG: fips140=only (strictest mode - panics on non-FIPS crypto)")
	} else if fipsMode == "on" {
		fmt.Println("GODEBUG: fips140=on (standard FIPS enforcement)")
	} else {
		fmt.Println("GODEBUG: Not set (default FIPS enforcement if built with GOFIPS140)")
	}

	fmt.Println("")
	fmt.Println("To enable FIPS compliance:")
	fmt.Println("  Build with: GOFIPS140=v1.0.0 go build")
	fmt.Println("  Run:        ./fips-client (FIPS automatically enforced)")
	fmt.Println()

	// Create a custom HTTP client with TLS configuration
	// In FIPS mode, only FIPS-approved cipher suites will be used
	// Note: FIPS 140-3 mode requires Extended Master Secret for TLS 1.2
	// If the server doesn't support EMS, the connection will fail
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // FIPS requires TLS 1.2 or higher
		// Let the FIPS library negotiate the best available version and cipher suite
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

	// Target URL
	url := "https://dibaddoo.snowflakecomputing.com/"
	fmt.Printf("\nSending GET request to: %s\n", url)

	// Send GET request
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)

		// Check if this is the expected FIPS 140-3 Extended Master Secret error
		if contains(err.Error(), "FIPS 140-3 requires the use of Extended Master Secret") {
			fmt.Println("\n📋 NOTE: This error is EXPECTED and CORRECT in FIPS 140-3 mode.")
			fmt.Println("It means:")
			fmt.Println("  ✓ FIPS enforcement is working correctly")
			fmt.Println("  ✓ The server doesn't support Extended Master Secret (EMS)")
			fmt.Println("  ✓ The connection was properly rejected per FIPS 140-3 requirements")
			fmt.Println("\nFIPS 140-3 compliance requires EMS for TLS 1.2 connections.")
			fmt.Println("To connect to this server:")
			fmt.Println("  1. Server must implement EMS support (RFC 7627)")
			fmt.Println("  2. Or run without FIPS mode (not recommended for FIPS-required environments)")
		}
		log.Fatal("")
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print response details
	fmt.Printf("\n--- Response Details ---\n")
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Protocol: %s\n", resp.Proto)

	// Print TLS connection info
	if resp.TLS != nil {
		fmt.Printf("\n--- TLS Information ---\n")
		fmt.Printf("Version: %s\n", getTLSVersion(resp.TLS.Version))
		fmt.Printf("Cipher Suite: %s\n", getCipherSuite(resp.TLS.CipherSuite))
		fmt.Printf("Server Name: %s\n", resp.TLS.ServerName)
		fmt.Printf("Negotiated Protocol: %s\n", resp.TLS.NegotiatedProtocol)
	}

	fmt.Printf("\n--- Response Body ---\n")
	fmt.Printf("%s\n", string(body))
}

// getFIPSMode checks the FIPS mode setting from GODEBUG
// Returns "on", "only", or empty string if disabled
func getFIPSMode(godebug string) string {
	if strings.Contains(godebug, "fips140=only") {
		return "only"
	}
	if strings.Contains(godebug, "fips140=on") {
		return "on"
	}
	return ""
}

// getTLSVersion converts TLS version constant to string
func getTLSVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

// getCipherSuite converts cipher suite constant to string
func getCipherSuite(suite uint16) string {
	suites := map[uint16]string{
		tls.TLS_RSA_WITH_RC4_128_SHA:                      "TLS_RSA_WITH_RC4_128_SHA",
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:                 "TLS_RSA_WITH_3DES_EDE_CBC_SHA",
		tls.TLS_RSA_WITH_AES_128_CBC_SHA:                  "TLS_RSA_WITH_AES_128_CBC_SHA",
		tls.TLS_RSA_WITH_AES_256_CBC_SHA:                  "TLS_RSA_WITH_AES_256_CBC_SHA",
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256:               "TLS_RSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256:               "TLS_RSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384:               "TLS_RSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:              "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:          "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:          "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:                "TLS_ECDHE_RSA_WITH_RC4_128_SHA",
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:           "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:            "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:            "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:       "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:         "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:         "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:       "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:         "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:       "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256:   "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256: "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
		tls.TLS_AES_128_GCM_SHA256:                        "TLS_AES_128_GCM_SHA256",
		tls.TLS_AES_256_GCM_SHA384:                        "TLS_AES_256_GCM_SHA384",
		tls.TLS_CHACHA20_POLY1305_SHA256:                  "TLS_CHACHA20_POLY1305_SHA256",
	}

	if name, ok := suites[suite]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (0x%04x)", suite)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Note: In Go 1.24+, FIPS mode uses the native Go Cryptographic Module
// Environment variables:
//   Build-time: GOFIPS140=v1.0.0 (or =latest) - selects FIPS module version
//   Runtime:    GODEBUG=fips140=on - enables FIPS mode
//               GODEBUG=fips140=only - strict enforcement (requires EMS)
// No need for BoringCrypto or GOEXPERIMENT=boringcrypto
