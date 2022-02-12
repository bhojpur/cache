package core

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

func TestCalculateSizeSum(t *testing.T) {
	sizeVector := []uint32{1, 2, 3, 4} // Sum = 10
	sum := calculateSizeSum(sizeVector)
	t.Logf("sum: %d\n", sum)

	sizeVector = []uint32{25, 13, 342, 92, 12, 456, 2} // Sum = 942
	sum = calculateSizeSum(sizeVector)
	t.Logf("sum: %d\n", sum)
}

func TestSplitBytes(t *testing.T) {
	bytes := []byte("Hello, world!")
	sizes := []uint32{3, 3, 3, 3, 1}

	splitBytes, err := SplitBytes(bytes, sizes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(splitBytes)
}
