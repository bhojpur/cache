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
)

// ringConsumer is the user-defined object responsible for receiving and
// processing items in batches when buffers are drained.
type ringConsumer interface {
	Push([]uint64) bool
}

// ringStripe is a singular ring buffer that is not concurrent safe.
type ringStripe struct {
	cons ringConsumer
	data []uint64
	capa int
}

func newRingStripe(cons ringConsumer, capa int64) *ringStripe {
	return &ringStripe{
		cons: cons,
		data: make([]uint64, 0, capa),
		capa: int(capa),
	}
}

// Push appends an item in the ring buffer and drains (copies items and
// sends to Consumer) if full.
func (s *ringStripe) Push(item uint64) {
	s.data = append(s.data, item)
	// Decide if the ring buffer should be drained.
	if len(s.data) >= s.capa {
		// Send elements to consumer and create a new ring stripe.
		if s.cons.Push(s.data) {
			s.data = make([]uint64, 0, s.capa)
		} else {
			s.data = s.data[:0]
		}
	}
}

// ringBuffer stores multiple buffers (stripes) and distributes Pushed items
// between them to lower contention.
//
// This implements the "batching" process described in the BP-Wrapper paper
// (section III part A).
type ringBuffer struct {
	pool *sync.Pool
}

// newRingBuffer returns a striped ring buffer. The Consumer in ringConfig will
// be called when individual stripes are full and need to drain their elements.
func newRingBuffer(cons ringConsumer, capa int64) *ringBuffer {
	// LOSSY buffers use a very simple sync.Pool for concurrently reusing
	// stripes. We do lose some stripes due to GC (unheld items in sync.Pool
	// are cleared), but the performance gains generally outweigh the small
	// percentage of elements lost. The performance primarily comes from
	// low-level runtime functions used in the standard library that aren't
	// available to us (such as runtime_procPin()).
	return &ringBuffer{
		pool: &sync.Pool{
			New: func() interface{} { return newRingStripe(cons, capa) },
		},
	}
}

// Push adds an element to one of the internal stripes and possibly drains if
// the stripe becomes full.
func (b *ringBuffer) Push(item uint64) {
	// Reuse or create a new stripe.
	stripe := b.pool.Get().(*ringStripe)
	stripe.Push(item)
	b.pool.Put(stripe)
}
