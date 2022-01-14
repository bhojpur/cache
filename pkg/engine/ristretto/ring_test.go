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
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConsumer struct {
	push func([]uint64)
	save bool
}

func (c *testConsumer) Push(items []uint64) bool {
	if c.save {
		c.push(items)
		return true
	}
	return false
}

func TestRingDrain(t *testing.T) {
	drains := 0
	r := newRingBuffer(&testConsumer{
		push: func(items []uint64) {
			drains++
		},
		save: true,
	}, 1)
	for i := 0; i < 100; i++ {
		r.Push(uint64(i))
	}
	require.Equal(t, 100, drains, "buffers shouldn't be dropped with BufferItems == 1")
}

func TestRingReset(t *testing.T) {
	drains := 0
	r := newRingBuffer(&testConsumer{
		push: func(items []uint64) {
			drains++
		},
		save: false,
	}, 4)
	for i := 0; i < 100; i++ {
		r.Push(uint64(i))
	}
	require.Equal(t, 0, drains, "testConsumer shouldn't be draining")
}

func TestRingConsumer(t *testing.T) {
	mu := &sync.Mutex{}
	drainItems := make(map[uint64]struct{})
	r := newRingBuffer(&testConsumer{
		push: func(items []uint64) {
			mu.Lock()
			defer mu.Unlock()
			for i := range items {
				drainItems[items[i]] = struct{}{}
			}
		},
		save: true,
	}, 4)
	for i := 0; i < 100; i++ {
		r.Push(uint64(i))
	}
	l := len(drainItems)
	require.NotEqual(t, 0, l)
	require.True(t, l <= 100)
}
