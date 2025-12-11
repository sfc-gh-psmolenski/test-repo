//go:build !boringcrypto

package main

func checkBoringCrypto() bool {
	return false
}

