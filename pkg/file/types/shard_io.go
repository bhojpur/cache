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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bhojpur/cache/pkg/file/common"
)

// ErrInvalidShard is returned when a shard's hash is not valid when loading from memory.
var ErrInvalidShard = errors.New("shard from memory is not valid")

// WriteShardToMemory writes a shard to memory.
func (shard *Shard) WriteShardToMemory() error {
	bytes := shard.Serialize()

	// Create a dir to store the shards
	err := common.CreateDirIfDoesNotExist("data/shards")
	if err != nil {
		return err
	}

	// Create the filename of the hash
	shardHashString := (*shard).Hash.String()[0:8]
	filename := fmt.Sprintf("data/shards/shard_%s.json", shardHashString)

	err = ioutil.WriteFile(filepath.FromSlash(filename), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReadShardFromMemory reads a shard from memory.
func ReadShardFromMemory(hash string) (*Shard, error) {
	// Read the file from memory
	data, err := ioutil.ReadFile(fmt.Sprintf("data/shards/shard_%s.json", hash))
	if err != nil {
		return nil, err
	}

	buffer := &Shard{} // Init a shard buffer

	// Read into the buffer
	err = json.Unmarshal(data, buffer)
	if err != nil {
		return nil, err
	}

	if (*buffer).Validate() == false {
		return nil, ErrInvalidShard
	}

	return buffer, nil // Return the shard pointer
}
