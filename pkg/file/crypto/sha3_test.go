package crypto

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

import "testing"

func TestNewHash(t *testing.T) {
	byteHash := []byte("11cd54753fc9e1d82e39f3b6f9727a3cc4cdf58eec127ccbe056829b1e0a9962")
	newHash := newHash(byteHash)

	t.Log(newHash)
}

func TestIsNil(t *testing.T) {
	byteHash := []byte("11cd54753fc9e1d82e39f3b6f9727a3cc4cdf58eec127ccbe056829b1e0a9962")
	hash := newHash(byteHash)

	nilByteHash := []byte("")
	nilHash := newHash(nilByteHash)

	if hash.IsNil() {
		t.Fatal("hash is not actually nil")
	}

	if nilHash.IsNil() == false {
		t.Fatal("hash is actually nil")
	}
}

func TestBytes(t *testing.T) {
	byteHash := []byte("11cd54753fc9e1d82e39f3b6f9727a3cc4cdf58eec127ccbe056829b1e0a9962")
	hash := newHash(byteHash)

	t.Log(hash.Bytes())
}

func TestString(t *testing.T) {
	byteHash := []byte("11cd54753fc9e1d82e39f3b6f9727a3cc4cdf58eec127ccbe056829b1e0a9962")
	hash := newHash(byteHash)

	t.Log(hash.String())
}
