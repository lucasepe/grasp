package csprng

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"math/rand"
	"strings"

	"golang.org/x/crypto/argon2"
)

// SecureSource contains a `cipher.Stream`.
type SecureSource struct{ cipher.Stream }

// Seed panics if invoked
func (s *SecureSource) Seed(seed int64) {
	panic("do not use")
}

// Int63 generates a (pseudo-)random 63bit value.
func (s *SecureSource) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Uint64 generates a (pseudo-)random 64bit value.
// Uses XORKeyStream with a zero input to extract
// the raw keystream as 64-bit integers.
func (s *SecureSource) Uint64() uint64 {
	var buf [8]byte
	s.XORKeyStream(buf[:], buf[:])
	return binary.LittleEndian.Uint64(buf[:])
}

// New returns a Source that can generate
// uniformly-distributed pseudo-random uint64 values.
//
// The first secret will be the Argon2 salt.
// Common use cases:
//    reddit.com my master password
//    example.com my master password
// here the site's name serves as the salt, which
// appropriately fits the role of the salt (a non-secret value).
// Argon2 is configured to output 32+16 bytes: 32 bytes for the
// key, 16 bytes for the IV. In this case (deterministic CSPRNG)
// the IV isn't important and it could simply be zero.
func New(secrets []string) (rand.Source, error) {
	const (
		memory      = 64 * 1024
		iterations  = 3
		parallelism = 2
		saltLength  = 16
		keyLength   = 32 + 16
	)
	salt := []byte(secrets[0])
	secret := []byte(strings.Join(secrets[1:], "\x00"))

	key := argon2.IDKey(secret, salt, iterations, memory, parallelism, keyLength)

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, key[32:])
	return &SecureSource{stream}, nil
}
