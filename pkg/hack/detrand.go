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
	_ "unsafe"
)

// DisableProtoBufRandomness disables the random insertion of whitespace characters when
// serializing Protocol Buffers in textual form (both when serializing to JSON or to ProtoText)
//
// Since the introduction of the APIv2 for Protocol Buffers, the default serializers in the
// package insert random whitespace characters that don't change the meaning of the serialized
// code but make byte-wise comparison impossible. The rationale behind this decision is as follows:
//
// "The ProtoBuf authors believe that golden tests are Wrong"
//
// Fine. Unfortunately, Vitess makes extensive use of golden tests through its test suite, which
// expect byte-wise comparison to be stable between test runs. Using the new version of the
// package would require us to rewrite hundreds of tests, or alternatively, we could disable
// the randomness and call it a day. The method required to disable the randomness is not public, but
// that won't stop us because we're good at computers.
//
//go:linkname DisableProtoBufRandomness google.golang.org/protobuf/internal/detrand.Disable
func DisableProtoBufRandomness()
