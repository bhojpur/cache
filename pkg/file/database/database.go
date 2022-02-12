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
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/bhojpur/cache/pkg/file/common"
	"github.com/bhojpur/cache/pkg/file/models"
	"github.com/bhojpur/cache/pkg/file/types"
	memcache "github.com/bhojpur/cache/pkg/memory"
)

// DBType represents the type of database of the instance, either a node's
// shard database or the main file database.
type DBType int

const (
	// FILEDB the marker for a file database.
	FILEDB DBType = iota

	// NSHARDDB is the marker for the node's shard database.
	NSHARDDB DBType = iota
)

// fileDBBucket is the bucket used to store files within a filedb.
var fileDBBucket = []byte("files")

// nshardDBBucket is the bucket used to store shards within a shard database.
var nshardDBBucket = []byte("shards")

// Database implements a general database that holds various data within meros.
type Database struct {
	Header types.DatabaseHeader `json:"header"` // Database header info
	Name   string               `json:"name"`   // The name of the db
	DBType DBType               // The type of database

	DB     *memcache.DB // Bhojpur Cache in-memory databse instance
	bucket []byte       // The bucket for the database
	open   bool         // Status of the DB
}

// Open opens the database for reading and writing. Creates a new DB if one
// with that name does not already exist.
func Open(dbName string, dbType DBType) (*Database, error) {
	// Make sure path exists
	err := common.CreateDirIfDoesNotExist(path.Join(models.DBPath, dbName))
	if err != nil {
		return nil, err
	}

	var database *Database // The database to return

	// Prepare to serialize the database struct
	databasePath := path.Join(models.DBPath, dbName, "db.json")
	if _, err := os.Stat(databasePath); err != nil { // If DB name does not exist
		// Create the database struct
		database = &Database{
			Header: types.NewDatabaseHeader(dbName), // Generate and set header

			Name: dbName, // Set the name
		}

		err = database.serialize(databasePath) // Write the db struct to disk
		if err != nil {
			return nil, err
		}

	} else {
		// If the db does exist, read from it and return it
		database, err = deserialize(databasePath)
		if err != nil {
			return nil, err
		}
	}

	// Prepare to open the Bhojpur Cache in-memory database
	memdbPath := path.Join(models.DBPath, dbName, "bhojpur-cache.db")
	db, err := memcache.Open(memdbPath, 0600, &memcache.Options{ // Open the DB
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	database.DB = db         // Set the DB
	database.DBType = dbType // Set the type of database

	// Bucket handler
	switch dbType {
	case FILEDB: // If FileDB type, set to the corresponding bucket
		database.bucket = fileDBBucket
	case NSHARDDB: // If NodeShardDB type, set to the corresponding bucket
		database.bucket = nshardDBBucket
	}

	err = database.makeBuckets() // Make the buckets in the database
	if err != nil {
		return nil, err
	}

	database.open = true // Set the status to open

	return database, nil
}

// Close closes the database.
func (db *Database) Close() error {
	err := db.DB.Close() // Close the DB
	if err != nil {
		return err
	}

	db.open = false // Set DB status
	return nil
}

// makeBuckets constructs the buckets in the database.
func (db *Database) makeBuckets() error {
	// Create all buckets in the database
	err := db.DB.Update(func(tx *memcache.Tx) error { // Open tx for bucket creation
		_, err := tx.CreateBucketIfNotExists(db.bucket) // Create bucket
		return err                                      // Handle err
	})
	if err != nil { // Check the err
		return err
	}
	return err
}

// String marshals the DB as a string.
func (db *Database) String() string {
	json, _ := json.MarshalIndent(*db, "", "  ")
	return string(json)
}

// serialize will serialize the database and write it to disk.
func (db *Database) serialize(filepath string) error {
	json, _ := json.MarshalIndent(*db, "", "  ")
	err := ioutil.WriteFile(filepath, json, 0600)
	return err
}

// deserialize will deserialize the database from the disk
func deserialize(filepath string) (*Database, error) {
	data, err := ioutil.ReadFile(filepath) // Read the database from disk
	if err != nil {
		return nil, err
	}

	buffer := &Database{} // Initialize the database buffer

	// Unmarshal and write into the buffer
	err = json.Unmarshal(data, buffer)

	return buffer, err
}
