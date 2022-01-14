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

// String force casts a []byte to a string.
// USE AT YOUR OWN RISK
func String(b []byte) (s string) {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// StringPointer returns &s[0], which is not allowed in go
func StringPointer(s string) unsafe.Pointer {
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(pstring.Data)
}

// StringBytes returns the underlying bytes for a string. Modifying this byte slice
// will lead to undefined behavior.
func StringBytes(s string) []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	hdr.Cap = len(s)
	hdr.Len = len(s)
	return b
}

// StringClone returns a newly allocated copy of the string that doesn't share
// its underlying memory storage.
func StringClone(s string) string {
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}
