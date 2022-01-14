package ristretto

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"math/rand"
	"time"
)

// cmSketch is a Count-Min sketch implementation with 4-bit counters
type cmSketch struct {
	rows [cmDepth]cmRow
	seed [cmDepth]uint64
	mask uint64
}

const (
	// cmDepth is the number of counter copies to store (think of it as rows).
	cmDepth = 4
)

func newCmSketch(numCounters int64) *cmSketch {
	if numCounters == 0 {
		panic("cmSketch: bad numCounters")
	}
	// Get the next power of 2 for better cache performance.
	numCounters = next2Power(numCounters)
	sketch := &cmSketch{mask: uint64(numCounters - 1)}
	// Initialize rows of counters and seeds.
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < cmDepth; i++ {
		sketch.seed[i] = source.Uint64()
		sketch.rows[i] = newCmRow(numCounters)
	}
	return sketch
}

// Increment increments the count(ers) for the specified key.
func (s *cmSketch) Increment(hashed uint64) {
	for i := range s.rows {
		s.rows[i].increment((hashed ^ s.seed[i]) & s.mask)
	}
}

// Estimate returns the value of the specified key.
func (s *cmSketch) Estimate(hashed uint64) int64 {
	min := byte(255)
	for i := range s.rows {
		val := s.rows[i].get((hashed ^ s.seed[i]) & s.mask)
		if val < min {
			min = val
		}
	}
	return int64(min)
}

// Reset halves all counter values.
func (s *cmSketch) Reset() {
	for _, r := range s.rows {
		r.reset()
	}
}

// Clear zeroes all counters.
func (s *cmSketch) Clear() {
	for _, r := range s.rows {
		r.clear()
	}
}

// cmRow is a row of bytes, with each byte holding two counters.
type cmRow []byte

func newCmRow(numCounters int64) cmRow {
	return make(cmRow, numCounters/2)
}

func (r cmRow) get(n uint64) byte {
	return byte(r[n/2]>>((n&1)*4)) & 0x0f
}

func (r cmRow) increment(n uint64) {
	// Index of the counter.
	i := n / 2
	// Shift distance (even 0, odd 4).
	s := (n & 1) * 4
	// Counter value.
	v := (r[i] >> s) & 0x0f
	// Only increment if not max value (overflow wrap is bad for LFU).
	if v < 15 {
		r[i] += 1 << s
	}
}

func (r cmRow) reset() {
	// Halve each counter.
	for i := range r {
		r[i] = (r[i] >> 1) & 0x77
	}
}

func (r cmRow) clear() {
	// Zero each counter.
	for i := range r {
		r[i] = 0
	}
}

func (r cmRow) string() string {
	s := ""
	for i := uint64(0); i < uint64(len(r)*2); i++ {
		s += fmt.Sprintf("%02d ", (r[(i/2)]>>((i&1)*4))&0x0f)
	}
	s = s[:len(s)-1]
	return s
}

// next2Power rounds x up to the next power of 2, if it's not already one.
func next2Power(x int64) int64 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}
