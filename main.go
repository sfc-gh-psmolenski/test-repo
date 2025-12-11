package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Check if FIPS mode is enabled
	if isFIPSEnabled() {
		fmt.Println("✓ FIPS mode is ENABLED")
	} else {
		fmt.Println("⚠ WARNING: FIPS mode is NOT enabled")
		fmt.Println("  To enable FIPS mode, build with: GOEXPERIMENT=boringcrypto go build")
	}

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
			fmt.Println("\n📋 NOTE: This error is EXPECTED when using strict FIPS 140-3 mode.")
			fmt.Println("It means:")
			fmt.Println("  ✓ FIPS enforcement is working correctly")
			fmt.Println("  ✓ The server doesn't support Extended Master Secret (EMS)")
			fmt.Println("  ✓ The connection was properly rejected per FIPS 140-3 requirements")
			fmt.Println("\nTo connect to this server, you would need to either:")
			fmt.Println("  1. Use a server that supports EMS (recommended)")
			fmt.Println("  2. Comment out 'crypto/tls/fipsonly' in boring.go for compatible FIPS mode")
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

// isFIPSEnabled checks if the Go runtime was built with FIPS mode
func isFIPSEnabled() bool {
	// Check if built with BoringCrypto support
	// This is determined at build time via build tags
	return checkBoringCrypto()
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
