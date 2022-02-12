package temporal

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
	"fmt"
	"os"
	"sync"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

type MemCacheConfig struct {
	Path string
	Mode os.FileMode
}

type TimeEntry struct {
	Time  int64
	Value []byte
}

type Query struct {
	Series string

	From int64
	To   int64
	// Sorting order:
	// Possible values are ASC and DESC
	// ASC : The time Series will have the oldest data first
	// DESC: The time Series will have the latest  data first.
	Sort string

	// Number of entries to be returned per page. This is used for pagination.
	// The next sequence is found out using NextEntry variable of a query response.
	MaxEntries int
}

type TimeSeries []TimeEntry

type Database interface {

	// Create a new bucket
	Create(name string) error

	// This function adds the records
	Add(name string, timeseries TimeSeries) error

	// Get the records
	Query(q Query) (timeSeries TimeSeries, nextEntry *int64, err error)

	QueryOnChannel(q Query) (timeseries <-chan TimeEntry, nextEntry chan *int64, err chan error)
	// Get the total pages for a particular query.
	// This helps for any client to call multiple queries
	GetPages(q Query) (seriesList []int64, count int, err error)

	// Get the records
	Get(series string) (timeSeries TimeSeries, err error)
	// Returns two channels, one for Time entries and one for error.
	// This avoids the usage of an extra buffer by the database
	// Caution: first read the channel and then read the error. Error channel shall be written only after the timeseries channel is closed
	GetOnChannel(series string) (timeseries <-chan TimeEntry, err chan error)

	// Delete a complete Series
	Delete(series string) error

	// Close the database
	Close() error
}

var ds Database    // It will be used as a singleton DB object
var once sync.Once // make thread safe singleton

func Open(config interface{}) (Database, error) {
	switch config.(type) {
	case MemCacheConfig:
		retDB := new(MemCacheDB)
		err := retDB.open(config.(MemCacheConfig))
		return retDB, err
	default:
		return nil, fmt.Errorf("Unsupported storage configuration")
	}
}
