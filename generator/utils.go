package generator

import (
	"math/rand"
)

// randomLetter extracts a random letter from the given string
func randomLetter(src *rand.Rand, s string) string {
	idx := src.Int63n(int64(len(s)))
	return string(s[idx])
}

// randomInsert randomly inserts the given value into the given string.
func randomInsert(src *rand.Rand, s, val string) string {
	if s == "" {
		return val
	}

	idx := src.Int63n(int64(len(s) + 1))

	return s[0:idx] + val + s[idx:]
}
