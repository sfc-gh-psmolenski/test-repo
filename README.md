# Go FIPS Mode HTTPS Client

This is a simple Go program that sends HTTPS GET requests with FIPS 140-3 mode enabled using Go's native FIPS support.

## What is FIPS Mode?

FIPS (Federal Information Processing Standards) 140-3 is a U.S. government security standard that specifies requirements for cryptographic modules. When FIPS mode is enabled, only FIPS-approved cryptographic algorithms and cipher suites are used.

## Go 1.24+ Native FIPS Support

Starting with Go 1.24, Go includes a **native FIPS 140-3 compliant cryptographic module** built directly into the standard library. This eliminates the need for BoringCrypto or external dependencies.

The Go Cryptographic Module v1.0.0 is currently undergoing NIST FIPS 140-3 validation and is listed in the CMVP Modules In Process List.

**References:**
- [Go Blog: FIPS 140-3](https://go.dev/blog/fips140)
- [Go Security: FIPS 140](https://tip.golang.org/doc/security/fips140)

## Requirements

- Go 1.25 (or Go 1.24+)
- `GODEBUG=fips140=on` environment variable to enable FIPS mode

## Building the Program

### Option 1: Using the build script (Recommended)

```bash
./build-fips.sh
```

### Option 2: Manual build

```bash
go build -o fips-client
```

**Note:** With Go 1.24+, FIPS mode is enabled at **runtime** via the `GODEBUG` environment variable, not at build time.

## Running the Program

### With FIPS mode enabled:

```bash
GODEBUG=fips140=on ./fips-client
```

Or set it permanently in your shell:

```bash
export GODEBUG=fips140=on
./fips-client
```

### Without FIPS mode (standard Go crypto):

```bash
./fips-client
```

The program will:
1. Check if FIPS mode is enabled (via GODEBUG environment variable)
2. Send a GET request to https://dibaddoo.snowflakecomputing.com/
3. Display TLS connection details including cipher suite used
4. Show the response from the server

## Verification

When running with `GODEBUG=fips140=on`, you should see:
- "✓ FIPS mode is ENABLED (Go native FIPS 140-3)" message
- TLS 1.2 or 1.3 connection
- FIPS-approved cipher suites (e.g., TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)

When running without the GODEBUG flag, you'll see:
- "⚠ WARNING: FIPS mode is NOT enabled" message
- Standard Go crypto libraries will be used

## FIPS-Approved Cipher Suites

In FIPS mode, only the following cipher suites are allowed:
- TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
- TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
- TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
- TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
- TLS_RSA_WITH_AES_128_GCM_SHA256
- TLS_RSA_WITH_AES_256_GCM_SHA384

## Troubleshooting

### FIPS mode is NOT enabled

If you see the warning message, it means the `GODEBUG=fips140=on` environment variable is not set. Enable it with:
```bash
GODEBUG=fips140=on ./fips-client
```

### Connection Errors

If you encounter connection errors, check:
1. Network connectivity
2. Firewall settings
3. TLS version compatibility (FIPS requires TLS 1.2+)

## Notes

- Go 1.24+ includes native FIPS 140-3 support (no BoringCrypto needed)
- FIPS mode is enabled at runtime via `GODEBUG=fips140=on` environment variable
- FIPS mode restricts the use of non-approved cryptographic algorithms
- Some cipher suites (like ChaCha20-Poly1305) are not FIPS-approved and won't be used in FIPS mode
- The program imports `crypto/tls/fipsonly` which enforces **strict FIPS 140-3 mode**
- FIPS 140-3 **requires Extended Master Secret (EMS)** for TLS 1.2 connections
- If a server doesn't support EMS, you'll get: `"tls: FIPS 140-3 requires the use of Extended Master Secret"` - **this is the correct behavior** and shows FIPS enforcement is working

## FIPS Mode Behavior in Go 1.25

Go 1.25's native FIPS 140-3 implementation enforces **strict compliance** by default:
- ✓ Only FIPS-approved cipher suites
- ✓ **Requires Extended Master Secret (EMS)** for TLS 1.2 connections
- ✓ Will reject non-compliant servers (expected behavior per FIPS 140-3)

This means that when you run with `GODEBUG=fips140=on`, connections to servers that don't support EMS (like some Snowflake instances) will fail with:
```
tls: FIPS 140-3 requires the use of Extended Master Secret
```

**This is correct behavior and demonstrates that FIPS enforcement is working properly.**

### GODEBUG Options

- `GODEBUG=fips140=on` - Enables FIPS mode with strict EMS enforcement
- Standard Go crypto (no GODEBUG flag) - No FIPS enforcement

## Implementation Details

### FIPS Mode Detection

The program checks the `GODEBUG` environment variable for `fips140=on`:

```go
func isFIPSEnabled() bool {
    godebug := os.Getenv("GODEBUG")
    return strings.Contains(godebug, "fips140=on")
}
```

### TLS Configuration

- Minimum TLS version: 1.2 (required by FIPS)
- Only FIPS-approved cipher suites are used when FIPS mode is enabled
- Example FIPS-approved cipher: `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`
- The `crypto/tls/fipsonly` import enforces strict compliance

## Output Examples

### When running with FIPS 140-3 mode enabled:

With a server that doesn't support Extended Master Secret (like Snowflake):

```bash
GODEBUG=fips140=on ./fips-client
```

Output:
```
✓ FIPS mode is ENABLED (standard mode: fips140=on)
  Note: Go 1.25 FIPS mode enforces Extended Master Secret for TLS 1.2

Sending GET request to: https://dibaddoo.snowflakecomputing.com/
Error making GET request: Get "https://dibaddoo.snowflakecomputing.com/": tls: FIPS 140-3 requires the use of Extended Master Secret

📋 NOTE: This error is EXPECTED and CORRECT in FIPS 140-3 mode.
It means:
  ✓ FIPS enforcement is working correctly
  ✓ The server doesn't support Extended Master Secret (EMS)
  ✓ The connection was properly rejected per FIPS 140-3 requirements

FIPS 140-3 compliance requires EMS for TLS 1.2 connections.
To connect to this server:
  1. Server must implement EMS support (RFC 7627)
  2. Or run without FIPS mode (not recommended for FIPS-required environments)
```

**This error is EXPECTED and CORRECT** - it demonstrates that FIPS 140-3 enforcement is working properly.

### When running without FIPS mode:

Standard Go crypto (no FIPS enforcement):

```bash
./fips-client
```

Output:
```
⚠ WARNING: FIPS mode is NOT enabled
  To enable FIPS mode:
    - GODEBUG=fips140=on ./fips-client

Sending GET request to: https://dibaddoo.snowflakecomputing.com/

--- Response Details ---
Status Code: 200
Status: 200 OK
Protocol: HTTP/1.1

--- TLS Information ---
Version: TLS 1.2
Cipher Suite: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
Server Name: dibaddoo.snowflakecomputing.com
Negotiated Protocol: 

--- Response Body ---
<html><head><title>Cookie support required</title>...
```

This uses standard Go cryptography without FIPS enforcement, allowing connections to any server.

