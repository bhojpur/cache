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
	"errors"

	"github.com/bhojpur/cache/pkg/file/crypto"
)

// errNilDBLabel is returned when a nil label is given.
var ErrNilDBLabel = errors.New("label for creating a shard database header must not be nil")

// shardDB is the database that holds the locations of each shard of a (larger) file.
type shardDB struct {
	header   DatabaseHeader        // Database header
	shardMap map[shardID]shardData // Shard data map

	hash crypto.Hash // Hash of the entire database
}

// generateShardDB constructs a new shard database.
func generateShardDB(shards []Shard, nodes []NodeID) (shardDB, error) {
	if len(shards) != len(nodes) {
		return shardDB{}, errors.New("shard count and node count do not match")
	}

	// Construct the map
	shardMap := make(map[shardID]shardData)

	// Generate and add shard data to the map
	for i, shard := range shards {
		id, data := generateShardEntry(shard, i, nodes[i]) // Generate the pair
		shardMap[id] = data                                // Put the data in the map
	}

	// Construct the database
	sharddb := shardDB{
		header:   NewDatabaseHeader(""), // Generate and set the header
		shardMap: shardMap,              // Set the shardMap
	}

	return sharddb, nil
}
