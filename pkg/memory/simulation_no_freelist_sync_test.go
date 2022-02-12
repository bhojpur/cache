package memory_test

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

	memcache "github.com/bhojpur/cache/pkg/memory"
)

func TestSimulateNoFreeListSync_1op_1p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 1, 1)
}
func TestSimulateNoFreeListSync_10op_1p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 10, 1)
}
func TestSimulateNoFreeListSync_100op_1p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 100, 1)
}
func TestSimulateNoFreeListSync_1000op_1p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 1000, 1)
}
func TestSimulateNoFreeListSync_10000op_1p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 10000, 1)
}
func TestSimulateNoFreeListSync_10op_10p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 10, 10)
}
func TestSimulateNoFreeListSync_100op_10p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 100, 10)
}
func TestSimulateNoFreeListSync_1000op_10p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 1000, 10)
}
func TestSimulateNoFreeListSync_10000op_10p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 10000, 10)
}
func TestSimulateNoFreeListSync_100op_100p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 100, 100)
}
func TestSimulateNoFreeListSync_1000op_100p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 1000, 100)
}
func TestSimulateNoFreeListSync_10000op_100p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 10000, 100)
}
func TestSimulateNoFreeListSync_10000op_1000p(t *testing.T) {
	testSimulate(t, &memcache.Options{NoFreelistSync: true}, 8, 10000, 1000)
}
