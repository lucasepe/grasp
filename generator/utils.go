package generator

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/lucasepe/grasp/mt19337"
	"golang.org/x/crypto/pbkdf2"
)

// mt19337WithKey initialize the Mersenne Twister PRNG
// using the specified slice of bytes.
func mt19337WithKey(key []byte) (*mt19337.MT19937, error) {
	var seed int64
	err := binary.Read(bytes.NewBuffer(key[:]), binary.BigEndian, &seed)
	if err != nil {
		return nil, err
	}

	rng := mt19337.New()
	rng.Seed(seed)

	return rng, nil
}

// deriveKey derives a key from []string slice.
// Returns a []byte of length 16 that can be used as cryptographic key.
// The key is derived based on the method described as PBKDF2 with the HMAC SHA256.
// The last item of the secrets array is the salt (the slice length must be greater then 1).
func deriveKey(secrets []string) ([]byte, error) {
	// len(secrets) must be > 1
	lastIdx := len(secrets) - 1
	if lastIdx <= 0 {
		return nil, fmt.Errorf("at least two words are required to initialize the generator")
	}

	// Generate the key bytes
	secret := strings.Join(secrets[:lastIdx], ".")
	salt := secrets[lastIdx]
	res := pbkdf2.Key([]byte(secret), []byte(salt), 2048, 16, sha256.New)

	return res, nil
}

// randomLetter extracts a random letter from the given string
func randomLetter(mt *mt19337.MT19937, s string) string {
	idx := mt.Int63n(int64(len(s)))
	return string(s[idx])
}

// randomInsert randomly inserts the given value into the given string.
func randomInsert(mt *mt19337.MT19937, s, val string) string {
	if s == "" {
		return val
	}

	idx := mt.Int63n(int64(len(s) + 1))

	return s[0:idx] + val + s[idx:]
}
