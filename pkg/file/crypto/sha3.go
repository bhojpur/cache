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

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// HashLength is the standardized length of a hash.
const HashLength = 32

// Hash represents the streamlined hash type to be used.
type Hash [HashLength]byte

// NewHash constructs a new hash given a hash, API so it returns an error.
func NewHash(b []byte) (Hash, error) {
	var hash Hash // Setup the hash
	bCropped := b // Setup the cropped buffer

	// Check the crop side
	if len(b) > len(hash) {
		bCropped = bCropped[len(bCropped)-HashLength:] // Crop the hash
	}

	// Copy the source
	copy(
		hash[HashLength-len(bCropped):],
		bCropped,
	)

	return hash, nil
}

// newHash constructs a new hash given a hash, returns no error
func newHash(b []byte) Hash {
	var hash Hash // Setup the hash
	bCropped := b // Setup the cropped buffer

	// Check the crop side
	if len(b) > len(hash) {
		bCropped = bCropped[len(bCropped)-HashLength:] // Crop the hash
	}

	// Copy the source
	copy(
		hash[HashLength-len(bCropped):],
		bCropped,
	)

	return hash
}

// Sha3 hashes a []byte using sha3.
func Sha3(b []byte) Hash {
	hash := sha3.New256()
	hash.Write(b)
	return newHash(hash.Sum(nil))
}

// Sha3String hashes a given message via sha3 and encodes the hashed message to a hex string.
func Sha3String(b []byte) string {
	b = Sha3(b).Bytes()
	return hex.EncodeToString(b) // Convert to a hex string
}

// HashFromString returns a Hash type given a hex string.
func HashFromString(s string) (Hash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, err
	}
	return newHash(b), nil
}

// IsNil checks if a given hash is nil.
func (hash Hash) IsNil() bool {
	nilBytes := 0 // Init nil bytes buffer

	// Iterate through the hash, checking for nil bytes
	for _, byteVal := range hash[:] {
		if byteVal == 0 {
			nilBytes++
		}
	}

	return nilBytes == HashLength
}

// Bytes converts a given hash to a byte array.
func (hash Hash) Bytes() []byte {
	return hash[:] // Return byte array value
}

// String returns the hash as a hex string.
func (hash Hash) String() string {
	b := hash.Bytes()
	return hex.EncodeToString(b) // Convert to a hex string
}
