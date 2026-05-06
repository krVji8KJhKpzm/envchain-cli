package main

import (
	"crypto/rand"
	"encoding/base64"
)

// generateToken returns a URL-safe base64-encoded random token
// of approximately the requested byte length.
// The actual string length may be slightly longer due to base64 encoding.
func generateToken(length int) string {
	if length <= 0 {
		length = 32
	}
	// We need ceil(length * 6/8) random bytes to produce `length` base64 chars.
	// Using length directly as byte count gives a slightly longer string; we trim.
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback: should never happen on a healthy system.
		panic("envchain: failed to read random bytes: " + err.Error())
	}
	encoded := base64.URLEncoding.EncodeToString(bytes)
	if len(encoded) > length {
		return encoded[:length]
	}
	return encoded
}
