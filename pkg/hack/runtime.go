package hack

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
	"reflect"
	"unsafe"
)

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

//go:noescape
//go:linkname strhash runtime.strhash
func strhash(p unsafe.Pointer, h uintptr) uintptr

// RuntimeMemhash provides access to the Go runtime's default hash function for arbitrary bytes.
// This is an optimal hash function which takes an input seed and is potentially implemented in hardware
// for most architectures. This is the same hash function that the language's `map` uses.
func RuntimeMemhash(b []byte, seed uint64) uint64 {
	pstring := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return uint64(memhash(unsafe.Pointer(pstring.Data), uintptr(seed), uintptr(pstring.Len)))
}

// RuntimeStrhash provides access to the Go runtime's default hash function for strings.
// This is an optimal hash function which takes an input seed and is potentially implemented in hardware
// for most architectures. This is the same hash function that the language's `map` uses.
func RuntimeStrhash(str string, seed uint64) uint64 {
	return uint64(strhash(unsafe.Pointer(&str), uintptr(seed)))
}

//go:linkname roundupsize runtime.roundupsize
func roundupsize(size uintptr) uintptr

// RuntimeAllocSize returns size of the memory block that mallocgc will allocate if you ask for the size.
func RuntimeAllocSize(size int64) int64 {
	return int64(roundupsize(uintptr(size)))
}

//go:linkname ParseFloatPrefix strconv.parseFloatPrefix
func ParseFloatPrefix(s string, bitSize int) (float64, int, error)
