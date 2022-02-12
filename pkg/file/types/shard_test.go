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
	"testing"

	"github.com/bhojpur/cache/pkg/file/crypto"
)

func TestNewShard(t *testing.T) {
	bytes := []byte("test bytes")
	newShard, err := NewShard(bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log((*newShard).String())
}
func TestCalculateShardSizes(t *testing.T) {
	rawBytes := []byte("123456789")
	nodes := 5

	sizes, _ := calculateShardSizes(rawBytes, nodes)
	t.Log(sizes)

}

func TestGenerateShards(t *testing.T) {
	bytes := []byte("these are the bytes of a test file that is going to be sharded.")
	nodes := 5

	shards, err := GenerateShards(
		bytes,
		nodes,
	)
	if err != nil {
		t.Fatal(err)
	}

	for i, shard := range shards {
		t.Logf("\n[shard %d] %s\n", i, shard.String())
	}

}

func TestFromBytes(t *testing.T) {
	bytes := []byte("test bytes")
	shard, err := NewShard(bytes)
	if err != nil {
		t.Fatal(err)
	}

	shardBytes := shard.Serialize()

	newShard, err := ShardFromBytes(shardBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("[newShard] %s\n", newShard.String())

}

func TestValidate(t *testing.T) {
	bytes := []byte("test bytes")
	shard, err := NewShard(bytes)
	if err != nil {
		t.Fatal(err)
	}

	isValid := (*shard).Validate()
	if isValid == false {
		t.Fatal("shard is actually valid")
	}

	badHash := crypto.Sha3([]byte("not a valid hash"))
	(*shard).Hash = badHash

	isValid = (*shard).Validate()
	if isValid == true {
		t.Fatal("shard is actually invalid")
	}
}
