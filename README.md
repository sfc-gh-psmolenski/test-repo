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
- `GOFIPS140=v1.0.0` environment variable at build time

**Important:** Building with `GOFIPS140` automatically enables FIPS 140-3 enforcement. The `GODEBUG=fips140` variable provides additional control but is optional.

## Building the Program

### Option 1: Using the build script (Recommended)

```bash
./build-fips.sh
```

This builds with `GOFIPS140=v1.0.0` by default.

### Option 2: Manual build with specific FIPS module version

```bash
# Build with FIPS module v1.0.0 (frozen in early 2025)
GOFIPS140=v1.0.0 go build -o fips-client

# Or use the latest FIPS module
GOFIPS140=latest go build -o fips-client
```

### Option 3: Build without FIPS module (standard crypto)

```bash
go build -o fips-client
```

**Important:** 
- `GOFIPS140` at **build time** enables FIPS 140-3 enforcement automatically
- When built with `GOFIPS140`, the program ALWAYS uses FIPS-approved cryptography
- `GODEBUG=fips140` at **runtime** provides additional control (optional)

## Running the Program

### Complete FIPS workflow (Build + Run):

```bash
# Build with FIPS module (FIPS is automatically enforced)
GOFIPS140=v1.0.0 go build -o fips-client

# Run - FIPS enforcement is already active from build
./fips-client
```

Or use the build script:

```bash
./build-fips.sh  # Builds with GOFIPS140=v1.0.0
./fips-client     # FIPS automatically enforced
```

### Optional: Runtime control with GODEBUG

```bash
# Standard FIPS enforcement (default when built with GOFIPS140)
./fips-client

# Explicit FIPS mode (same as above if built with GOFIPS140)
GODEBUG=fips140=on ./fips-client

# Strictest mode (panics on any non-FIPS crypto)
GODEBUG=fips140=only ./fips-client
```

### Without FIPS (build without GOFIPS140):

```bash
# Build without FIPS module
go build -o fips-client

# Run with standard Go crypto (no FIPS enforcement)
./fips-client
```

The program will:
1. Display the FIPS module version used at build time (GOFIPS140)
2. Check if FIPS mode is enabled at runtime (GODEBUG)
3. Send a GET request to https://dibaddoo.snowflakecomputing.com/
4. Display TLS connection details including cipher suite used
5. Show the response from the server

## Verification

### Build Verification
The program shows which FIPS module was used at build time:
- `ℹ️  Built with FIPS module: GOFIPS140=v1.0.0` - Built with FIPS support
- `ℹ️  Built without GOFIPS140 set` - Built with standard crypto

### Runtime Verification
When running with `GODEBUG=fips140=on`, you should see:
- "✓ FIPS mode is ENABLED (standard mode: GODEBUG=fips140=on)" message
- TLS 1.2 or 1.3 connection
- FIPS-approved cipher suites (e.g., TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)

When running without the GODEBUG flag, you'll see:
- "⚠ WARNING: FIPS runtime mode is NOT enabled" message
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

If you see the warning message, it means either:
1. The program wasn't built with `GOFIPS140` set, or
2. The `GODEBUG=fips140=on` runtime variable is not set

To enable both:
```bash
# Build with FIPS module
GOFIPS140=v1.0.0 go build -o fips-client

# Run with FIPS mode enabled
GODEBUG=fips140=on ./fips-client
```

### Connection Errors

If you encounter connection errors, check:
1. Network connectivity
2. Firewall settings
3. TLS version compatibility (FIPS requires TLS 1.2+)

## Notes

- Go 1.24+ includes native FIPS 140-3 support (no BoringCrypto needed)
- **GOFIPS140 at build time is the primary control:**
  - Building with `GOFIPS140=v1.0.0` **automatically enables** FIPS 140-3 enforcement
  - The built binary will ALWAYS use FIPS-approved cryptography
  - `GODEBUG=fips140` provides optional additional runtime control
- FIPS mode restricts the use of non-approved cryptographic algorithms
- Some cipher suites (like ChaCha20-Poly1305) are not FIPS-approved and won't be used in FIPS mode
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

### Environment Variables

**Build Time (Primary Control):**
- `GOFIPS140=v1.0.0` - Build with FIPS module v1.0.0 → **FIPS always enforced**
- `GOFIPS140=latest` - Build with latest FIPS module → **FIPS always enforced**
- No GOFIPS140 - Build with standard Go crypto → **No FIPS enforcement**

**Runtime (Optional Additional Control):**
- No GODEBUG - Uses default FIPS enforcement (if built with GOFIPS140)
- `GODEBUG=fips140=on` - Explicit FIPS mode (same as default if built with GOFIPS140)
- `GODEBUG=fips140=only` - Strictest mode (panics on any non-FIPS crypto)

## Implementation Details

### FIPS Mode Detection

The program checks both build-time and runtime environment variables:

```go
// Build-time: Check FIPS module version
gofips140 := os.Getenv("GOFIPS140")

// Runtime: Check FIPS mode enabled
godebug := os.Getenv("GODEBUG")
fipsEnabled := strings.Contains(godebug, "fips140=on")
```

### TLS Configuration

- Minimum TLS version: 1.2 (required by FIPS)
- Only FIPS-approved cipher suites are used when FIPS mode is enabled
- Example FIPS-approved cipher: `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`
- Extended Master Secret (EMS) is required for TLS 1.2 in FIPS mode

## Output Examples

### When built with FIPS (GOFIPS140=v1.0.0):

With a server that doesn't support Extended Master Secret (like Snowflake):

```bash
# Build with FIPS
GOFIPS140=v1.0.0 go build -o fips-client

# Run (FIPS automatically enforced)
./fips-client
```

Output:
```
=== FIPS Mode Status ===

ℹ️  FIPS enforcement is determined at BUILD time with GOFIPS140=v1.0.0
ℹ️  When built with GOFIPS140, FIPS 140-3 compliance is ALWAYS enforced

GODEBUG: Not set (default FIPS enforcement if built with GOFIPS140)

Sending GET request to: https://dibaddoo.snowflakecomputing.com/
Error making GET request: Get "https://dibaddoo.snowflakecomputing.com/": tls: FIPS 140-3 requires the use of Extended Master Secret

📋 NOTE: This error is EXPECTED and CORRECT in FIPS 140-3 mode.
It means:
  ✓ FIPS enforcement is working correctly
  ✓ The server doesn't support Extended Master Secret (EMS)
  ✓ The connection was properly rejected per FIPS 140-3 requirements
```

**This error is EXPECTED and CORRECT** - it demonstrates that FIPS 140-3 enforcement is working properly.

### When built WITHOUT FIPS (no GOFIPS140):

Standard Go crypto (no FIPS enforcement):

```bash
# Build without FIPS
go build -o fips-client

# Run
./fips-client
```

Output:
```
=== FIPS Mode Status ===

ℹ️  FIPS enforcement is determined at BUILD time with GOFIPS140=v1.0.0
ℹ️  When built with GOFIPS140, FIPS 140-3 compliance is ALWAYS enforced

GODEBUG: Not set (default FIPS enforcement if built with GOFIPS140)

To enable FIPS compliance:
  Build with: GOFIPS140=v1.0.0 go build
  Run:        ./fips-client (FIPS automatically enforced)

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

