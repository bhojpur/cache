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
	"errors"

	"github.com/bhojpur/cache/pkg/file/types"
	memcache "github.com/bhojpur/cache/pkg/memory"
)

// generateEntry generates an ID-file/shard pair for the DB.
func generateEntry(item interface{}) (ID, []byte, error) {
	// Type assertion (disambiguation)
	if t, ok := item.(types.File); ok {
		return ID(t.Hash), t.Bytes(), nil
	} else if t, ok := item.(types.Shard); ok {
		return ID(t.Hash), t.Serialize(), nil
	} else {
		return ID{}, nil, errors.New("invalid type to store in database")
	}
}

// PutItem adds a new item to the database.
func (db *Database) PutItem(item interface{}) (ID, error) {
	var t ID              // Temporary nil item ID
	if db.open == false { // Make sure the DB is open
		return t, errors.New("database is closed")
	}

	// Extract the data for the database
	id, data, err := generateEntry(item)
	if err != nil {
		return ID{}, err
	}

	// Write the item to the bucket
	err = db.DB.Update(func(tx *memcache.Tx) error {
		b := tx.Bucket(db.bucket) // Fetch the bucket

		// Put necessary data into the bucket
		return b.Put(id.Bytes(), data)
	})

	return id, err
}

// GetItem gets an item from the database.
func (db *Database) GetItem(id ID) (interface{}, error) {
	if db.open == false { // Make sure the DB is open
		return nil, errors.New("db is closed")
	}

	// Initialize buffer
	var buffer []byte

	// Read from the database
	err := db.DB.View(func(tx *memcache.Tx) error {
		b := tx.Bucket(db.bucket)   // Fetch the bucket
		dbRead := b.Get(id.Bytes()) // Read the item from the db
		if dbRead == nil {          // Check the item not nil
			return errors.New(
				"item '" + id.String() + "' not found in db '" + db.Name + "'",
			) // Return err if nil
		}

		buffer = make([]byte, len(dbRead)) // Init the buffer size
		copy(buffer, dbRead)               // Copy the item to the buffer
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Construct corresponding type from bytes and return
	switch db.DBType {
	case FILEDB:
		return types.FileFromBytes(buffer)
	case NSHARDDB:
		return types.ShardFromBytes(buffer)
	}

	// Throw for undefined behavior
	return nil, errors.New("invalid database type was used")
}
