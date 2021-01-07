// mt19937.go - an implementation of the 64bit Mersenne Twister PRNG
// Copyright (C) 2013  Jochen Voss <voss@seehuhn.de>
// Changed on 2021 Luca Sepe <luca.sepe@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mt19337

const (
	n         = 312
	m         = 156
	notSeeded = n + 1

	hiMask uint64 = 0xffffffff80000000
	loMask uint64 = 0x000000007fffffff

	matrixA uint64 = 0xB5026F5AA96619E9
)

// MT19937 is the structure to hold the state of one instance of the
// Mersenne Twister PRNG.  New instances can be allocated using the
// mt19937.New() function.  MT19937 implements the rand.Source
// interface and rand.New() from the math/rand package can be used to
// generate different distributions from a MT19937 PRNG.
//
// This class is not safe for concurrent accesss by different
// goroutines.  If more than one goroutine accesses the PRNG, the
// callers must synchronise access using sync.Mutex or similar.
type MT19937 struct {
	state []uint64
	index int
}

// New allocates a new instance of the 64bit Mersenne Twister.
// A seed can be set using the .Seed() or .SeedFromSlice() methods.
// If no seed is set explicitly, a default seed is used instead.
func New() *MT19937 {
	res := &MT19937{
		state: make([]uint64, n),
		index: notSeeded,
	}
	return res
}

// Seed uses the given 64bit value to initialise the generator state.
// This method is part of the rand.Source interface.
func (mt *MT19937) Seed(seed int64) {
	x := mt.state
	x[0] = uint64(seed)
	for i := uint64(1); i < n; i++ {
		x[i] = 6364136223846793005*(x[i-1]^(x[i-1]>>62)) + i
	}
	mt.index = n
}

// SeedFromSlice uses the given slice of 64bit values to set the
// generator state.
func (mt *MT19937) SeedFromSlice(key []uint64) {
	mt.Seed(19650218)

	x := mt.state
	i := uint64(1)
	j := 0
	k := len(key)
	if n > k {
		k = n
	}
	for k > 0 {
		x[i] = x[i] ^ ((x[i-1] ^ (x[i-1] >> 62)) * 3935559000370003845) + key[j] + uint64(j)
		i++
		if i >= n {
			x[0] = x[n-1]
			i = 1
		}
		j++
		if j >= len(key) {
			j = 0
		}
		k--
	}
	for j := uint64(0); j < n-1; j++ {
		x[i] = x[i] ^ ((x[i-1] ^ (x[i-1] >> 62)) * 2862933555777941757) - i
		i++
		if i >= n {
			x[0] = x[n-1]
			i = 1
		}
	}
	x[0] = 1 << 63
}

// Uint64 generates a (pseudo-)random 64bit value.  The output can be
// used as a replacement for a sequence of independent, uniformly
// distributed samples in the range 0, 1, ..., 2^64-1.  This method is
// part of the rand.Source64 interface.
func (mt *MT19937) Uint64() uint64 {
	x := mt.state
	if mt.index >= n {
		if mt.index == notSeeded {
			mt.Seed(5489) // default seed, as in mt19937-64.c
		}
		for i := 0; i < n-m; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		y := (x[n-1] & hiMask) | (x[0] & loMask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixA)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return y
}

// Int63 generates a (pseudo-)random 63bit value.  The output can be
// used as a replacement for a sequence of independent, uniformly
// distributed samples in the range 0, 1, ..., 2^63-1.  This method is
// part of the rand.Source interface.
func (mt *MT19937) Int63() int64 {
	x := mt.state
	if mt.index >= n {
		if mt.index == notSeeded {
			mt.Seed(5489) // default seed, as in mt19937-64.c
		}
		for i := 0; i < n-m; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		y := (x[n-1] & hiMask) | (x[0] & loMask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixA)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= (y >> 29) & 0x5555555555555555
	y ^= (y << 17) & 0x71D67FFFEDA60000
	y ^= (y << 37) & 0xFFF7EEE000000000
	y ^= (y >> 43)
	mt.index++
	return int64(y & 0x7fffffffffffffff)
}

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32.
func (mt *MT19937) Int31() int32 { return int32(mt.Int63() >> 32) }

// Int31n returns, as an int32, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func (mt *MT19937) Int31n(n int32) int32 {
	if n <= 0 {
		panic("invalid argument to Int31n")
	}

	if n&(n-1) == 0 { // n is power of two, can mask
		return mt.Int31() & (n - 1)
	}

	max := int32((1 << 31) - 1 - (1<<31)%uint32(n))
	v := mt.Int31()

	for v > max {
		v = mt.Int31()
	}

	return v % n
}

// int31n returns, as an int32, a non-negative pseudo-random number in [0,n).
// n must be > 0, but int31n does not check this; the caller must ensure it.
// int31n exists because Int31n is inefficient, but Go 1 compatibility
// requires that the stream of values produced by math/rand remain unchanged.
// int31n can thus only be used internally, by newly introduced APIs.
//
// For implementation details, see:
// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction
// https://lemire.me/blog/2016/06/30/fast-random-shuffling
func (mt *MT19937) int31n(n int32) int32 {
	v := mt.Uint32()

	prod := uint64(v) * uint64(n)
	low := uint32(prod)

	if low < uint32(n) {
		thresh := uint32(-n) % uint32(n)

		for low < thresh {
			v = mt.Uint32()
			prod = uint64(v) * uint64(n)
			low = uint32(prod)
		}
	}

	return int32(prod >> 32)
}

// Uint32 returns a pseudo-random 32-bit value as a uint32.
func (mt *MT19937) Uint32() uint32 { return uint32(mt.Int63() >> 31) }

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func (mt *MT19937) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(mt.Int63n(int64(i + 1)))
		swap(i, j)
	}

	for ; i > 0; i-- {
		j := int(mt.int31n(int32(i + 1)))
		swap(i, j)
	}
}

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
func (mt *MT19937) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")

	}

	if n&(n-1) == 0 { // n is power of two, can mask
		return mt.Int63() & (n - 1)
	}

	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))

	v := mt.Int63()

	for v > max {
		v = mt.Int63()
	}

	return v % n
}

// Read fills `p` with (pseudo-)random bytes.  This method implements
// the io.Reader interface.  The returned length `n` always equals
// `len(p)` and `err` is always nil.
func (mt *MT19937) Read(p []byte) (n int, err error) {
	n = len(p)
	for len(p) >= 8 {
		val := mt.Uint64()
		p[0] = byte(val)
		p[1] = byte(val >> 8)
		p[2] = byte(val >> 16)
		p[3] = byte(val >> 24)
		p[4] = byte(val >> 32)
		p[5] = byte(val >> 40)
		p[6] = byte(val >> 48)
		p[7] = byte(val >> 56)
		p = p[8:]
	}
	if len(p) > 0 {
		val := mt.Uint64()
		for i := 0; i < len(p); i++ {
			p[i] = byte(val)
			val >>= 8
		}
	}
	return n, nil
}
