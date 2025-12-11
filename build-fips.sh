#!/bin/bash

# Build script for FIPS-enabled Go client
# This script builds the program with BoringCrypto support for FIPS compliance

echo "Building Go FIPS client..."
GOEXPERIMENT=boringcrypto go build -o fips-client

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "Run the program with: ./fips-client"
else
    echo "✗ Build failed"
    exit 1
fi

