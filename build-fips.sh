#!/bin/bash

# Build script for FIPS-enabled Go client
# Go 1.24+ uses native FIPS 140-3 support (no BoringCrypto needed)

echo "Building Go FIPS client with native FIPS 140-3 support..."
go build -o fips-client

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "To run with FIPS mode enabled:"
    echo "  GODEBUG=fips140=on ./fips-client"
    echo ""
    echo "Or set it permanently:"
    echo "  export GODEBUG=fips140=on"
    echo "  ./fips-client"
else
    echo "✗ Build failed"
    exit 1
fi

