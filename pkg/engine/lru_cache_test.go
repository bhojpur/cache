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
	"testing"
)

type CacheValue struct {
	size int64
}

func cacheValueSize(val interface{}) int64 {
	return val.(*CacheValue).size
}

func TestInitialState(t *testing.T) {
	cache := NewLRUCache(5, cacheValueSize)
	l, sz, c, e := cache.Len(), cache.UsedCapacity(), cache.MaxCapacity(), cache.Evictions()
	if l != 0 {
		t.Errorf("length = %v, want 0", l)
	}
	if sz != 0 {
		t.Errorf("size = %v, want 0", sz)
	}
	if c != 5 {
		t.Errorf("capacity = %v, want 5", c)
	}
	if e != 0 {
		t.Errorf("evictions = %v, want 0", c)
	}
}

func TestSetInsertsValue(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	data := &CacheValue{0}
	key := "key"
	cache.Set(key, data)

	v, ok := cache.Get(key)
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value: %v != %v", data, v)
	}

	values := cache.Items()
	if len(values) != 1 || values[0].Key != key {
		t.Errorf("Cache.Values() returned incorrect values: %v", values)
	}
}

func TestGetValueWithMultipleTypes(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	data := &CacheValue{0}
	key := "key"
	cache.Set(key, data)

	v, ok := cache.Get("key")
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value for \"key\": %v != %v", data, v)
	}

	v, ok = cache.Get(string([]byte{'k', 'e', 'y'}))
	if !ok || v.(*CacheValue) != data {
		t.Errorf("Cache has incorrect value for []byte {'k','e','y'}: %v != %v", data, v)
	}
}

func TestSetUpdatesSize(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	emptyValue := &CacheValue{0}
	key := "key1"
	cache.Set(key, emptyValue)
	if sz := cache.UsedCapacity(); sz != 0 {
		t.Errorf("cache.UsedCapacity() = %v, expected 0", sz)
	}
	someValue := &CacheValue{20}
	key = "key2"
	cache.Set(key, someValue)
	if sz := cache.UsedCapacity(); sz != 20 {
		t.Errorf("cache.UsedCapacity() = %v, expected 20", sz)
	}
}

func TestSetWithOldKeyUpdatesValue(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	emptyValue := &CacheValue{0}
	key := "key1"
	cache.Set(key, emptyValue)
	someValue := &CacheValue{20}
	cache.Set(key, someValue)

	v, ok := cache.Get(key)
	if !ok || v.(*CacheValue) != someValue {
		t.Errorf("Cache has incorrect value: %v != %v", someValue, v)
	}
}

func TestSetWithOldKeyUpdatesSize(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	emptyValue := &CacheValue{0}
	key := "key1"
	cache.Set(key, emptyValue)

	if sz := cache.UsedCapacity(); sz != 0 {
		t.Errorf("cache.UsedCapacity() = %v, expected %v", sz, 0)
	}

	someValue := &CacheValue{20}
	cache.Set(key, someValue)
	expected := int64(someValue.size)
	if sz := cache.UsedCapacity(); sz != expected {
		t.Errorf("cache.UsedCapacity() = %v, expected %v", sz, expected)
	}
}

func TestGetNonExistent(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)

	if _, ok := cache.Get("notthere"); ok {
		t.Error("Cache returned a notthere value after no inserts.")
	}
}

func TestDelete(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	value := &CacheValue{1}
	key := "key"

	cache.Delete(key)
	cache.Set(key, value)
	cache.Delete(key)

	if sz := cache.UsedCapacity(); sz != 0 {
		t.Errorf("cache.UsedCapacity() = %v, expected 0", sz)
	}

	if _, ok := cache.Get(key); ok {
		t.Error("Cache returned a value after deletion.")
	}
}

func TestClear(t *testing.T) {
	cache := NewLRUCache(100, cacheValueSize)
	value := &CacheValue{1}
	key := "key"

	cache.Set(key, value)
	cache.Clear()

	if sz := cache.UsedCapacity(); sz != 0 {
		t.Errorf("cache.UsedCapacity() = %v, expected 0 after Clear()", sz)
	}
}

func TestCapacityIsObeyed(t *testing.T) {
	size := int64(3)
	cache := NewLRUCache(100, cacheValueSize)
	cache.SetCapacity(size)
	value := &CacheValue{1}

	// Insert up to the cache's capacity.
	cache.Set("key1", value)
	cache.Set("key2", value)
	cache.Set("key3", value)
	if sz := cache.UsedCapacity(); sz != size {
		t.Errorf("cache.UsedCapacity() = %v, expected %v", sz, size)
	}
	// Insert one more; something should be evicted to make room.
	cache.Set("key4", value)
	sz, evictions := cache.UsedCapacity(), cache.Evictions()
	if sz != size {
		t.Errorf("post-evict cache.UsedCapacity() = %v, expected %v", sz, size)
	}
	if evictions != 1 {
		t.Errorf("post-evict cache.Evictions() = %v, expected 1", evictions)
	}

	// Check various other stats
	if l := cache.Len(); int64(l) != size {
		t.Errorf("cache.Len() returned bad length: %v", l)
	}
	if s := cache.UsedCapacity(); s != size {
		t.Errorf("cache.UsedCapacity() returned bad size: %v", s)
	}
	if c := cache.MaxCapacity(); c != size {
		t.Errorf("cache.UsedCapacity() returned bad length: %v", c)
	}
}

func TestLRUIsEvicted(t *testing.T) {
	size := int64(3)
	cache := NewLRUCache(size, cacheValueSize)

	cache.Set("key1", &CacheValue{1})
	cache.Set("key2", &CacheValue{1})
	cache.Set("key3", &CacheValue{1})
	// lru: [key3, key2, key1]

	// Look up the elements. This will rearrange the LRU ordering.
	cache.Get("key3")
	cache.Get("key2")
	cache.Get("key1")
	// lru: [key1, key2, key3]

	cache.Set("key0", &CacheValue{1})
	// lru: [key0, key1, key2]

	// The least recently used one should have been evicted.
	if _, ok := cache.Get("key3"); ok {
		t.Error("Least recently used element was not evicted.")
	}

	if e, want := cache.Evictions(), int64(1); e != want {
		t.Errorf("evictions: %d, want: %d", e, want)
	}
}
