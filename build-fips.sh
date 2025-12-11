#!/bin/bash

# Build script for FIPS-enabled Go client
# Go 1.24+ uses native FIPS 140-3 support (no BoringCrypto needed)

# Set FIPS module version (v1.0.0 or latest)
FIPS_VERSION="${GOFIPS140:-v1.0.0}"

echo "Building Go FIPS client with FIPS module version: $FIPS_VERSION"
echo ""

GOFIPS140=$FIPS_VERSION go build -o fips-client

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Build successful with GOFIPS140=$FIPS_VERSION"
    echo ""
    echo "FIPS 140-3 enforcement is NOW ACTIVE in this binary."
    echo ""
    echo "To run:"
    echo "  ./fips-client"
    echo ""
    echo "Optional - for strictest mode (panics on non-FIPS crypto):"
    echo "  GODEBUG=fips140=only ./fips-client"
else
    echo "✗ Build failed"
    exit 1
fi

