# Go FIPS Mode HTTPS Client

This is a simple Go program that sends HTTPS GET requests with FIPS mode enabled.

## What is FIPS Mode?

FIPS (Federal Information Processing Standards) 140-2 is a U.S. government security standard that specifies requirements for cryptographic modules. When FIPS mode is enabled, only FIPS-approved cryptographic algorithms and cipher suites are used.

## Requirements

- Go 1.25 or later
- GOEXPERIMENT=boringcrypto flag for FIPS mode

## Building with FIPS Mode

### Option 1: Using the build script (Recommended)

```bash
./build-fips.sh
```

### Option 2: Manual build with GOEXPERIMENT

```bash
GOEXPERIMENT=boringcrypto go build -o fips-client
```

### Option 3: Regular build (without FIPS)

```bash
go build -o fips-client
```

Note: When built without FIPS mode, the program will display a warning but will still function using standard Go crypto libraries.

## Running the Program

```bash
./fips-client
```

The program will:
1. Check if FIPS mode is enabled
2. Send a GET request to https://dibaddoo.snowflakecomputing.com/
3. Display TLS connection details including cipher suite used
4. Show the response from the server

## Verification

When running with FIPS mode enabled, you should see:
- "✓ FIPS mode is ENABLED" message
- TLS 1.2 or 1.3 connection
- FIPS-approved cipher suites (e.g., TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)

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

If you see the warning message, it means the program was built without BoringCrypto. Try rebuilding with:
```bash
GOEXPERIMENT=boringcrypto go build
```

### Connection Errors

If you encounter connection errors, check:
1. Network connectivity
2. Firewall settings
3. TLS version compatibility (FIPS requires TLS 1.2+)

## Notes

- BoringCrypto is Google's FIPS 140-2 validated fork of OpenSSL
- FIPS mode restricts the use of non-approved cryptographic algorithms
- Some cipher suites (like ChaCha20-Poly1305) are not FIPS-approved and won't be used in FIPS mode
- The program uses **strict FIPS 140-3 mode** with `crypto/tls/fipsonly` import
- FIPS 140-3 **requires Extended Master Secret (EMS)** for TLS 1.2 connections
- If a server doesn't support EMS, you'll get: `"tls: FIPS 140-3 requires the use of Extended Master Secret"` - **this is the correct behavior** and shows FIPS enforcement is working

## FIPS Modes

### Strict FIPS 140-3 Mode (Current Configuration)

The program is configured with `crypto/tls/fipsonly` import, which enforces the strictest FIPS 140-3 requirements:
- ✓ Only FIPS-approved cipher suites
- ✓ Requires Extended Master Secret for TLS 1.2
- ✓ Will reject non-compliant servers (expected behavior)

### Compatible FIPS Mode (Optional)

To allow connections to servers without EMS support while still using FIPS-validated crypto:
1. Comment out `import _ "crypto/tls/fipsonly"` in `boring.go`
2. Rebuild with `./build-fips.sh`

Note: This still uses BoringCrypto and FIPS-approved algorithms, but doesn't enforce EMS requirement.

## Implementation Details

### FIPS Mode Detection

The program uses Go build tags to detect if it was compiled with BoringCrypto:
- `boring.go`: Compiled only when `boringcrypto` build tag is present
- `noboring.go`: Compiled when `boringcrypto` build tag is absent

### TLS Configuration

- Minimum TLS version: 1.2 (required by FIPS)
- Only FIPS-approved cipher suites are used when built with BoringCrypto
- Example FIPS-approved cipher: `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`

## Output Examples

### When running with strict FIPS 140-3 mode:

With a server that doesn't support Extended Master Secret (like Snowflake):

```
✓ FIPS mode is ENABLED

Sending GET request to: https://dibaddoo.snowflakecomputing.com/
Error making GET request: Get "https://dibaddoo.snowflakecomputing.com/": tls: FIPS 140-3 requires the use of Extended Master Secret

📋 NOTE: This error is EXPECTED when using strict FIPS 140-3 mode.
It means:
  ✓ FIPS enforcement is working correctly
  ✓ The server doesn't support Extended Master Secret (EMS)
  ✓ The connection was properly rejected per FIPS 140-3 requirements

To connect to this server, you would need to either:
  1. Use a server that supports EMS (recommended)
  2. Comment out 'crypto/tls/fipsonly' in boring.go for compatible FIPS mode
```

**This error is EXPECTED and CORRECT** - it demonstrates that FIPS 140-3 enforcement is working properly. The program now clearly explains what the error means and why it occurs.

### When running with compatible FIPS mode:

After commenting out `crypto/tls/fipsonly` in `boring.go`:

```
✓ FIPS mode is ENABLED

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

This still uses FIPS-validated BoringCrypto but allows connections to servers without EMS.

