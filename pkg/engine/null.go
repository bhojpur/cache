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

// nullCache is a no-op cache that does not store items
type nullCache struct{}

// Get never returns anything on the nullCache
func (n *nullCache) Get(_ string) (interface{}, bool) {
	return nil, false
}

// Set is a no-op in the nullCache
func (n *nullCache) Set(_ string, _ interface{}) bool {
	return false
}

// ForEach iterates the nullCache, which is always empty
func (n *nullCache) ForEach(_ func(interface{}) bool) {}

// Delete is a no-op in the nullCache
func (n *nullCache) Delete(_ string) {}

// Clear is a no-op in the nullCache
func (n *nullCache) Clear() {}

// Wait is a no-op in the nullcache
func (n *nullCache) Wait() {}

func (n *nullCache) Len() int {
	return 0
}

// Capacity returns the capacity of the nullCache, which is always 0
func (n *nullCache) UsedCapacity() int64 {
	return 0
}

// Capacity returns the capacity of the nullCache, which is always 0
func (n *nullCache) MaxCapacity() int64 {
	return 0
}

// SetCapacity sets the capacity of the null cache, which is a no-op
func (n *nullCache) SetCapacity(_ int64) {}

func (n *nullCache) Evictions() int64 {
	return 0
}
