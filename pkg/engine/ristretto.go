package engine

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
	"github.com/bhojpur/cache/pkg/engine/ristretto"
)

var _ Cache = &ristretto.Cache{}

// NewRistrettoCache returns a Cache implementation based on Ristretto
func NewRistrettoCache(maxEntries, maxCost int64, cost func(interface{}) int64) *ristretto.Cache {
	// The TinyLFU paper recommends to allocate 10x times the max entries amount as counters
	// for the admission policy; since our caches are small and we're very interested on admission
	// accuracy, we're a bit more greedy than 10x
	const CounterRatio = 12

	config := ristretto.Config{
		NumCounters: maxEntries * CounterRatio,
		MaxCost:     maxCost,
		BufferItems: 64,
		Metrics:     true,
		Cost:        cost,
	}
	cache, err := ristretto.NewCache(&config)
	if err != nil {
		panic(err)
	}
	return cache
}
