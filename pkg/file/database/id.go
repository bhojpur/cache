package database

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

	"github.com/bhojpur/cache/pkg/file/crypto"
)

// ID represents a hash for the keys in the database.
type ID crypto.Hash

// IDFromString returns an ID given a string
func IDFromString(s string) (ID, error) {
	b, err := hex.DecodeString(s) // Decode from hex into []byte
	if err != nil {
		return ID{}, err
	}

	idHash, err := crypto.NewHash(b) // Create the hash
	return ID(idHash), err           // Return the cast to ID
}

// Bytes converts a given hash to a byte array.
func (id ID) Bytes() []byte {
	hash := crypto.Hash(id)
	return hash.Bytes() // Return byte array value
}

// String returns the hash as a hex string.
func (id ID) String() string {
	b := id.Bytes()
	return hex.EncodeToString(b) // Convert to a hex string
}
