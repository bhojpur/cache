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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/cache/pkg/engine/ristretto"
)

func TestNewDefaultCacheImpl(t *testing.T) {
	assertNullCache := func(t *testing.T, cache Cache) {
		_, ok := cache.(*nullCache)
		require.True(t, ok)
	}

	assertLFUCache := func(t *testing.T, cache Cache) {
		_, ok := cache.(*ristretto.Cache)
		require.True(t, ok)
	}

	assertLRUCache := func(t *testing.T, cache Cache) {
		_, ok := cache.(*LRUCache)
		require.True(t, ok)
	}

	tests := []struct {
		cfg    *Config
		verify func(t *testing.T, cache Cache)
	}{
		{&Config{MaxEntries: 0, MaxMemoryUsage: 0, LFU: false}, assertNullCache},
		{&Config{MaxEntries: 0, MaxMemoryUsage: 0, LFU: true}, assertNullCache},
		{&Config{MaxEntries: 100, MaxMemoryUsage: 0, LFU: false}, assertLRUCache},
		{&Config{MaxEntries: 0, MaxMemoryUsage: 1000, LFU: false}, assertNullCache},
		{&Config{MaxEntries: 100, MaxMemoryUsage: 1000, LFU: false}, assertLRUCache},
		{&Config{MaxEntries: 100, MaxMemoryUsage: 0, LFU: true}, assertNullCache},
		{&Config{MaxEntries: 100, MaxMemoryUsage: 1000, LFU: true}, assertLFUCache},
		{&Config{MaxEntries: 0, MaxMemoryUsage: 1000, LFU: true}, assertNullCache},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d.%d.%v", tt.cfg.MaxEntries, tt.cfg.MaxMemoryUsage, tt.cfg.LFU), func(t *testing.T) {
			cache := NewDefaultCacheImpl(tt.cfg)
			tt.verify(t, cache)
		})
	}
}
