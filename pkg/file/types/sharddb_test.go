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
	"math/rand"
	"testing"

	"github.com/bhojpur/cache/pkg/file/crypto"
)

func TestGenerateShardDB(t *testing.T) {
	// Generate shards
	testBytes := []byte("im a neat test file")
	shards, err := GenerateShards(testBytes, 3)
	if err != nil {
		t.Fatal(err)
	}

	// Make len(shards) random nodes
	var nodes []NodeID
	for i := 0; i < len(shards); i++ {
		randomPort := rand.Intn(10000-9000) + 9000
		randomNode, err := NewNodeID("0.0.0.0", randomPort)
		if err != nil {
			t.Fatal(err)
		}

		nodes = append(nodes, *randomNode)
	}

	// Generate shardDB
	sharddb, err := generateShardDB(shards, nodes)
	if err != nil {
		t.Fatal(err)
	}

	// Test the map
	for k, v := range sharddb.shardMap {
		t.Logf("%x: %v\n\n", k, v)
	}

	// Test map access
	h, err := crypto.HashFromString("bafc4c93a862aecc87368f291090fc2fe479eecc3bd15b7efdde01ff92c42592")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\nGET: [%v]\n", sharddb.shardMap[shardID(h)])

	// Analyze shard sizes
	var sizeCounter uint32
	for _, v := range sharddb.shardMap {
		t.Logf("size: %d\n", v.Size)
		sizeCounter = sizeCounter + v.Size
	}
	t.Logf("correct size: %d", len(testBytes))
	t.Logf(" actual size: %d", sizeCounter)
}
