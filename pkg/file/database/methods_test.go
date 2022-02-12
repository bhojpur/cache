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
	"testing"

	"github.com/bhojpur/cache/pkg/file/types"
)

func TestPutFile(t *testing.T) {
	filedb, err := Open("myfiledb", FILEDB)
	if err != nil {
		t.Fatal(err)
	}
	defer filedb.Close()

	file, err := types.NewFile("test_file.txt")
	if err != nil {
		t.Fatal(err)
	}

	fid, err := filedb.PutItem(*file)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fileID: %s\n", fid.String())
}

func TestGetFile(t *testing.T) {

	/* -- PUT -- */

	filedb, err := Open("myfiledb", FILEDB)
	if err != nil {
		t.Fatal(err)
	}
	defer filedb.Close()

	file, err := types.NewFile("test_file.txt")
	if err != nil {
		t.Fatal(err)
	}

	fid, err := filedb.PutItem(*file)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fileID: %s\n", fid.String())

	/* -- GET -- */

	readFile, err := filedb.GetItem(fid)
	if err != nil {
		t.Logf("%v", file)
		t.Fatal(err)
	}

	v := readFile.(*types.File)
	t.Logf("file '%s': %v\n", fid.String(), v)
}

func TestPutShard(t *testing.T) {
	sharddb, err := Open("mynodesharddb", NSHARDDB)
	if err != nil {
		t.Fatal(err)
	}
	defer sharddb.Close()

	shard, err := types.NewShard([]byte("I represent some bytes in a shard"))
	if err != nil {
		t.Fatal(err)
	}

	sid, err := sharddb.PutItem(*shard)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("shardID: %s\n", sid.String())
}

func TestGetShard(t *testing.T) {
	sharddb, err := Open("mynodesharddb", NSHARDDB)
	if err != nil {
		t.Fatal(err)
	}
	defer sharddb.Close()

	sid, err := IDFromString("f910670bbb0012b2eb6c4f321a02251d2f38c97a8c28acdda16bb0a3b79c1ab5")
	if err != nil {
		t.Fatal(err)
	}

	shard, err := sharddb.GetItem(sid)
	if err != nil {
		t.Fatal(err)
	}

	v := shard.(*types.Shard)
	t.Logf("shard '%s': %v\n", sid, v)
}
