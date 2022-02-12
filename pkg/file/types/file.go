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
	"errors"
	"os"

	"github.com/bhojpur/cache/pkg/file/crypto"
	"github.com/bhojpur/cache/pkg/file/models"
)

var (
	// ErrNilFilename is returned when the fileame to construct a new file is nil.
	ErrNilFilename = errors.New("filename to construct file must not be nil")

	// ErrNilShardCount is returned when the shard counnt to cosntruct a new file is nil.
	ErrNilShardCount = errors.New("shard count to cosntruct file must not be nil")

	// ErrNilFileSize is returned when the file size to construct a new file is nil.
	ErrNilFileSize = errors.New("file size to construct file must not be nil")
)

// File contains the (important) metadata of a file stored in a database.
type File struct {
	Filename   string      `json:"filename"`    // The file's filename
	ShardCount int         `json:"shard_count"` // The number of shards hosting the file
	Size       uint32      `json:"size"`        // Total size of the file
	ShardDB    *shardDB    `json:"shard_db"`    // Pointer to this file's shardDb
	Hash       crypto.Hash `json:"hash"`        // The hash of the file
}

// NewFile constructs a new file from a file in memory.
func NewFile(filename string) (*File, error) {
	// Check that the filename is not nil
	if filename == "" {
		return nil, ErrNilFilename
	}

	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the file length and bytes
	fileStat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := uint32(fileStat.Size())
	if size == 0 {
		return nil, ErrNilFileSize
	}

	// Read from the file
	bytes := make([]byte, size)
	_, err = f.Read(bytes)
	if err != nil {
		return nil, err
	}

	shardCount := models.ShardCount

	// Create a new file pointer
	file := &File{
		Filename:   filename,   // The filename
		ShardCount: shardCount, // The total amount of shards hostinng the file
		Size:       size,       // The total size of the file
		ShardDB:    nil,        // nil for now
	}

	// Compute the hash of the file
	(*file).Hash = crypto.Sha3(file.Bytes())
	return file, nil
}

// FileFromBytes constructs a *File from a []byte.
func FileFromBytes(b []byte) (*File, error) {
	if b == nil {
		return nil, errors.New("cannot construct file from nil []byte")
	}
	buffer := &File{}                // Init buffer
	err := json.Unmarshal(b, buffer) // Unmarshal json
	return buffer, err
}

/* ----- BEGIN HELPER FUNCTIONS ----- */

// Bytes converts the database header to bytes.
func (file *File) Bytes() []byte {
	json, _ := json.MarshalIndent(*file, "", "  ")
	return json
}

// String converts the database to a string.
func (file *File) String() string {
	json, _ := json.MarshalIndent(*file, "", "  ")
	return string(json)
}

/* ----- END HELPER FUNCTIONS ----- */
