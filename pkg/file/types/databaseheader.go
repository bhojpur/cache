package types

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
	"encoding/json"
	"time"

	"github.com/bhojpur/cache/pkg/file/crypto"
)

// DatabaseHeader is the identifier for a database.
type DatabaseHeader struct {
	Label   string      `json:"label"`   // A database label
	Created string      `json:"created"` // The time that the database was created
	Hash    crypto.Hash `json:"hash"`    // The hash of the database header
}

// NewDatabaseHeader creates a new database header.
func NewDatabaseHeader(label string) DatabaseHeader {
	// Create the header
	newDatabaseHeader := DatabaseHeader{
		Label:   label,
		Created: time.Now().String(), // The timestamp
	}

	// Compute the header hash and return
	newDatabaseHeader.Hash = crypto.Sha3(newDatabaseHeader.Bytes())
	return newDatabaseHeader
}

/* ----- BEGIN HELPER FUNCTIONS ----- */

// Bytes converts the database header to bytes.
func (databaseHeader DatabaseHeader) Bytes() []byte {
	json, _ := json.MarshalIndent(databaseHeader, "", "  ")
	return json
}

// String converts the database to a string.
func (databaseHeader DatabaseHeader) String() string {
	json, _ := json.MarshalIndent(databaseHeader, "", "  ")
	return string(json)
}

/* ----- END HELPER FUNCTIONS ----- */
