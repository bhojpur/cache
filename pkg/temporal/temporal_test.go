package temporal_test

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
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	temporal "github.com/bhojpur/cache/pkg/temporal"
)

func setupDB(dbName string) (temporal.Database, string, error) {
	temp_file := fmt.Sprintf("test_temp_%s", dbName)
	config := temporal.MemCacheConfig{Path: temp_file}

	db, err := temporal.Open(config)
	if err != nil {
		return nil, temp_file, err
	}
	return db, temp_file, nil
}

func clean(db temporal.Database, temp_filepath string) {
	db.Close()
	err := os.Remove(temp_filepath)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestTsdb_Add(t *testing.T) {
	tname := "TestTsdb_Add"
	db, filePath, err := setupDB(tname)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer clean(db, filePath)
	series := createDummyRecords(100, false)
	err = db.Add(tname, series)

	if err != nil {
		t.Fatal(err)
	}

	query := temporal.Query{Series: tname, MaxEntries: 100, From: series[0].Time, To: series[len(series)-1].Time, Sort: temporal.ASC}
	resSeries, nextEntry, err := db.Query(query)
	if err != nil {
		t.Fatal(err)
	}

	if compareTimeSeries(resSeries, series) == false {
		t.Error("Inserted and Fetched did not match")
	}
	if nextEntry != nil {
		t.Error("nextEntry is null")
	}
}

func TestMemCacheDB_Create(t *testing.T) {
	tname := "TestMemCacheDB_Create"
	db, filePath, err := setupDB(tname)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer clean(db, filePath)

	err = db.Create(tname)
	if err != nil {
		t.Error(err.Error())
	}

	err = db.Create(tname)
	if err == nil {
		t.Error("Duplicate Buckets are allowed")
	}
}
func TestTsdb_CheckDescending(t *testing.T) {
	tname := "TestTsdb_CheckDescending"
	db, filePath, err := setupDB(tname)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer clean(db, filePath)
	series := createDummyRecords(100, false)
	err = db.Add(tname, series)

	if err != nil {
		t.Fatal(err)
	}

	query := temporal.Query{Series: tname, MaxEntries: 100, From: series[0].Time, To: series[len(series)-1].Time, Sort: temporal.DESC}
	resSeries, nextEntry, err := db.Query(query)
	if err != nil {
		t.Fatal(err)
	}

	// reverse the Series
	for i := len(series)/2 - 1; i >= 0; i-- {
		opp := len(series) - 1 - i
		series[i], series[opp] = series[opp], series[i]
	}
	if compareTimeSeries(resSeries, series) == false {
		t.Error("Inserted and Fetched did not match")
	}
	if nextEntry != nil {
		t.Error("nextEntry is null")
	}
}

func TestTsdb_QueryPagination(t *testing.T) {
	tname := "TestTsdb_CheckDescending"
	db, filePath, err := setupDB(tname)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer clean(db, filePath)
	series := createDummyRecords(100, false)
	err = db.Add(tname, series)

	query := temporal.Query{Series: tname, MaxEntries: 50, From: series[0].Time, To: series[len(series)-1].Time, Sort: temporal.ASC}
	resSeries, nextEntry, err := db.Query(query)
	if err != nil {
		t.Fatal(err)
	}

	if compareTimeSeries(resSeries, series[0:50]) == false {
		t.Error("First page entries did not match")
	}

	if nextEntry == nil {
		t.Error("nextEntry is null")
	}
	query = temporal.Query{Series: tname, MaxEntries: 50, From: *nextEntry, To: series[len(series)-1].Time, Sort: temporal.ASC}
	resSeries, nextEntry, err = db.Query(query)
	if err != nil {
		t.Fatal(err)
	}

	if compareTimeSeries(resSeries, series[50:100]) == false {
		t.Error("Second page entries did not match")
	}

	if nextEntry != nil {
		t.Error("nextEntry is null")
	}
}

func TestTsdb_QueryPages(t *testing.T) {
	tname := "TestTsdb_CheckDescending"
	limit := 25
	count := 100
	db, filePath, err := setupDB(tname)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer clean(db, filePath)
	series := createDummyRecords(count, false)
	err = db.Add(tname, series)

	query := temporal.Query{Series: tname, MaxEntries: limit, From: series[0].Time, To: series[len(series)-1].Time, Sort: temporal.ASC}
	pages, retCount, err := db.GetPages(query)
	if err != nil {
		t.Fatal(err)
	}

	if retCount != count {
		t.Error("Length of time series do not match")
	}

	if len(pages) != count/limit {
		t.Error("Number of pages is not as expected")
	}

	for i := 0; i < count/limit; i = i + 1 {
		if pages[i] != series[i*limit].Time {
			t.Errorf("Page indices are not matching %v != %v", pages[i], series[i*limit].Time)
		}
	}
}

func compareTimeSeries(s1, s2 temporal.TimeSeries) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		e1 := s1[i]
		e2 := s2[i]
		if e1.Time != e2.Time || bytes.Compare(e1.Value, e2.Value) != 0 {
			return false
		}
	}

	return true
}

func createDummyRecords(count int, dec bool) temporal.TimeSeries {
	timeVal := time.Date(
		2009, 0, 0, 0, 0, 0, 0, time.UTC)

	timeseries := make(temporal.TimeSeries, 0, count)

	for i := 0; i < count; i++ {
		value, err := timeVal.MarshalBinary()
		if err != nil {
			log.Fatal(err)
			return nil
		}
		timeseries = append(timeseries, temporal.TimeEntry{timeVal.UnixNano(), value})
		if dec {
			timeVal = timeVal.Add(-time.Second)
		} else {
			timeVal = timeVal.Add(time.Second)
		}
	}
	return timeseries
}
