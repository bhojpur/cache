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
	"math"
	"time"

	"github.com/bhojpur/cache/pkg/file/core"
	"github.com/bhojpur/cache/pkg/file/crypto"
)

var (
	// ErrNilBytes is thrown when a shard is constructed when given nil bytes.
	ErrNilBytes = errors.New("bytes to construct new shard must not be nil")

	// ErrCannotCalculateShardSizes is thrown when the []byte to a CalculateShardSizes call is nil.
	ErrCannotCalculateShardSizes = errors.New("bytes to calculate shard sizes must not be nil")
)

// Shard is a struct that holds a piece of data that is
// a part of another, bigger piece of data.
type Shard struct {
	Size      uint32      `json:"size"`      // The size of the shard
	Bytes     []byte      `json:"bytes"`     // The actual data of the shard
	Hash      crypto.Hash `json:"hash"`      // The hash of the shard
	Timestamp string      `json:"timestamp"` // The timestamp of the shard
}

// NewShard attempts to construct a new shard.
func NewShard(bytes []byte) (*Shard, error) {
	if bytes == nil {
		return nil, ErrNilBytes
	}

	// Make the new shard
	newShard := &Shard{
		Size:      uint32(len(bytes)),
		Bytes:     bytes,
		Hash:      crypto.Sha3(bytes), // Hash the bytes of the shard, not the shard itself
		Timestamp: time.Now().String(),
	}

	return newShard, nil
}

// GenerateShards generates a slice of shards given a string of bytes
// This is an interface/api function.
func GenerateShards(bytes []byte, n int) ([]Shard, error) {
	var shards []Shard // Init the shard slice

	shardSizes, err := calculateShardSizes(bytes, n) // Calculate the shard sizes
	if err != nil {
		return nil, err
	}

	splitBytes, err := core.SplitBytes(bytes, shardSizes) // Split the bytes into the correct sizes
	if err != nil {
		return nil, err
	}

	// Generate the slices
	for i := 0; i < len(shardSizes); i++ {
		// Create a new shard
		newShard, err := NewShard(
			splitBytes[i],
		)
		if err != nil { // Check error
			return nil, err
		}
		shards = append(shards, *newShard) // Append the new shard to the shard slice
	}

	return shards, nil
}

// calculateShardSizes determines the recommended size of each shard.
func calculateShardSizes(raw []byte, n int) ([]uint32, error) {
	rawSize := len(raw)

	// Check that the input is not null
	if rawSize == 0 || n == 0 {
		return nil, ErrCannotCalculateShardSizes
	}

	partition := math.Floor(float64(rawSize / n)) // Calculate the size of each shard
	partitionSize := uint32(partition)            // Convert to a uint32
	modulo := uint32(rawSize % n)                 // Calculate the module mod n

	// Populate a slice of the correct shard sizes
	var sizes []uint32
	for i := 0; i < n; i++ {
		sizes = append(sizes, partitionSize)
	}

	// Adjust for the left over bytes
	if modulo+partitionSize >= partitionSize*uint32(n) {
		// This will be optimized eventually
	}

	sizes[n-1] += modulo // Add the left over bytes to the last element

	return sizes, nil
}

/* ----- BEGIN HELPER FUNCTIONS ----- */

func (shard *Shard) String() string {
	json, _ := json.MarshalIndent(*shard, "", "  ")
	return string(json)
}

// Serialize serializes a Shard pointer to bytes.
func (shard *Shard) Serialize() []byte {
	json, _ := json.MarshalIndent(*shard, "", "  ")
	return json
}

// Validate makes sure that the shard is valid.
func (shard *Shard) Validate() bool {
	if crypto.Sha3((*shard).Bytes) == (*shard).Hash {
		return true
	}
	return false
}

// ShardFromBytes constructs a *Shard from bytes.
func ShardFromBytes(b []byte) (*Shard, error) {
	buffer := &Shard{}               // Init shard buffer
	err := json.Unmarshal(b, buffer) // Unmarshal json
	return buffer, err
}

/* ----- END HELPER FUNCTIONS ----- */
