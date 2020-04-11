package utilities

import (
	"crypto/sha256"
)

//FastSha256 fast inline sha256()
func FastSha256(params ...[]byte) []byte {
	h := sha256.New()
	for i := 0; i < len(params); i++ {
		h.Write(params[i])
	}
	return h.Sum(nil)
}

//ByteArrayEqual is two byte array equal
func ByteArrayEqual(a []byte, b []byte) bool {
	if l := len(a); l == len(b) {
		c := byte(0)
		for i := 0; c == 0 && i < l; i++ {
			c = a[i] ^ b[i]
		}
		return c == 0
	}
	return false
}

//ByteArraySafeEqual timing safe equal
func ByteArraySafeEqual(a []byte, b []byte) bool {
	if l := len(a); l == len(b) {
		c := byte(0)
		for i := 0; i < l; i++ {
			c |= a[i] ^ b[i]
		}
		return c == 0
	}
	return false
}
