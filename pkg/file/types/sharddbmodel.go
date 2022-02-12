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
	"github.com/bhojpur/cache/pkg/file/crypto"
)

// shardID is the key model in the shard map (shard db).
type shardID crypto.Hash

// shardData is the value model in the shard map (shard db).
type shardData struct {
	Index  int    `json:"index"`   // The position of the shard out of all the shards for this file
	NodeID NodeID `json:"node_id"` // The ID of the node holding the shard
	Size   uint32 `json:"size"`    // The size in bytes of the shard
}

// generateShardEntry generates a shardID-shardData pair.
func generateShardEntry(shard Shard, index int, nodeID NodeID) (shardID, shardData) {
	shardData := shardData{
		Index:  index,      // Set the index
		NodeID: nodeID,     // Set the NodeID
		Size:   shard.Size, // Get and set the size
	}

	return shardID(shard.Hash), shardData // Return the pair
}
