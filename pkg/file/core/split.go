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

import "errors"

// ErrCannotSplitBytes is returned when a byte slice cannot be split evenly
// with the given size array.
var ErrCannotSplitBytes = errors.New(
	"byte array can not be split given size vector",
)

// calculateSizeSum calculates the sum of all of the uint32s in a []uint32.
func calculateSizeSum(sizes []uint32) uint32 {
	var sum uint32 // Init the sum

	// Add all of the slice's elements
	for _, size := range sizes {
		sum += size
	}

	return sum // Return the sum
}

// SplitBytes splits a []byte n times.
func SplitBytes(bytes []byte, sizes []uint32) ([][]byte, error) {
	// Check that the bytes can be split given the size vector
	if uint32(len(bytes)) != calculateSizeSum(sizes) {
		return nil, ErrCannotSplitBytes
	}

	var splitBytes [][]byte // Init the master slice
	currentBytePos := 0     // Init the byte position

	// For each size (shard)
	for _, currentSize := range sizes {
		var tempBytes []byte // Init shard[i]'s byte slice

		// For each byte that needs to be added
		for i := 0; i < int(currentSize); i++ {
			tempBytes = append(tempBytes, bytes[currentBytePos]) // Add the byte

			currentBytePos++ // Move the "byte cursor"
		}

		// Append the shard's bytes to the master slice
		splitBytes = append(splitBytes, tempBytes)
	}

	return splitBytes, nil
}
