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

import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	memcache "github.com/bhojpur/cache/pkg/memory"
)

type MemCacheDB struct {
	db *memcache.DB
}

func (mem *MemCacheDB) open(config MemCacheConfig) error {
	var err error
	mode := config.Mode
	if mode == 0 {
		mode = 0600
	}
	mem.db, err = memcache.Open(config.Path, mode, nil)
	return err
}

func (mem MemCacheDB) Close() error {
	return mem.db.Close()
}

// This function converts a floating point number to a bytearray
func timeToByteArr(timeVal int64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(timeVal))

	return buff

}

// This function converts a bytearray floating point number
func byteArrToTime(byteArr []byte) int64 {
	//This is set to bigendian so that the timestamp is sorted in binary format.
	timeVal := int64(binary.BigEndian.Uint64(byteArr))
	return timeVal
}

func (mem MemCacheDB) Create(name string) error {
	if name == "" {
		return fmt.Errorf("time Series record with Empty name")
	}
	return mem.db.Update(func(tx *memcache.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
func (mem MemCacheDB) Add(name string, timeseries TimeSeries) error {
	if name == "" {
		return fmt.Errorf("Time Series record with empty name")
	}
	if err := mem.db.Batch(func(tx *memcache.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		for _, entry := range timeseries {

			err = b.Put(timeToByteArr(entry.Time), entry.Value)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (mem MemCacheDB) Query(q Query) (timeSeries TimeSeries, nextEntry *int64, err error) {
	timeSeries = make([]TimeEntry, 0, q.MaxEntries)

	nextEntry = nil
	err = mem.db.View(func(tx *memcache.Tx) error {

		b := tx.Bucket([]byte(q.Series))
		if b == nil {
			return fmt.Errorf("Bucket:%v does not exist", q.Series)
		}

		c := b.Cursor()

		//Default case : If the sorting is descending
		first := q.To
		last := q.From
		next := c.Prev
		loopCondition := func(val int64, last int64) bool {
			return val >= last
		}
		//else
		if strings.Compare(q.Sort, ASC) == 0 {
			first = q.From
			last = q.To
			next = c.Next
			loopCondition = func(val int64, last int64) bool {
				return val <= last
			}

		}

		count := 0
		// Iterate over the time values
		var k, v []byte
		for k, v = c.Seek(timeToByteArr(first)); k != nil && loopCondition(byteArrToTime(k), last) && count < q.MaxEntries; k, v = next() {
			record := TimeEntry{byteArrToTime(k), v}
			timeSeries = append(timeSeries, record)
			count = count + 1
		}
		if count == q.MaxEntries && k != nil && loopCondition(byteArrToTime(k), last) {
			ne := byteArrToTime(k)
			nextEntry = &ne
		}
		return nil
	})

	if err != nil {
		return TimeSeries{}, nil, err
	}

	return timeSeries, nextEntry, nil
}

func (mem MemCacheDB) QueryOnChannel(q Query) (<-chan TimeEntry, chan *int64, chan error) {
	resultCh := make(chan TimeEntry, 10)
	errorCh := make(chan error)
	nextEntryChan := make(chan *int64)

	go func() {
		var nextEntry *int64
		err := mem.db.View(func(tx *memcache.Tx) error {

			b := tx.Bucket([]byte(q.Series))
			if b == nil {
				return fmt.Errorf("Bucket:%v does not exist", q.Series)
			}

			c := b.Cursor()
			count := 0
			if q.Sort == DESC {
				k, v := c.Seek(timeToByteArr(q.To))
				if k == nil { //if the seek value is beyond the last entry then go to the last entry
					k, v = c.Last()
				}

				start := timeToByteArr(q.From)
				// Iterate over the time values
				for ; k != nil && bytes.Compare(k, start) >= 0 && count != q.MaxEntries; k, v = c.Prev() {
					record := TimeEntry{byteArrToTime(k), v}
					resultCh <- record
					count++
				}
				if count == q.MaxEntries && k != nil && bytes.Compare(k, start) >= 0 {
					ne := byteArrToTime(k)
					nextEntry = &ne
				}
			} else {
				k, v := c.Seek(timeToByteArr(q.From))
				last := timeToByteArr(q.To)
				// Iterate over the time values
				for ; k != nil && bytes.Compare(k, last) <= 0 && count != q.MaxEntries; k, v = c.Next() {
					record := TimeEntry{byteArrToTime(k), v}
					resultCh <- record
					count = count + 1
				}
				if count == q.MaxEntries && k != nil && bytes.Compare(k, last) <= 0 {
					ne := byteArrToTime(k)
					nextEntry = &ne
				}
			}

			return nil
		})

		//make sure you close the resultchannel before error channel
		close(resultCh)
		nextEntryChan <- nextEntry

		if err != nil {
			errorCh <- err
		}
		close(errorCh)
	}()

	return resultCh, nextEntryChan, errorCh
}

func (mem MemCacheDB) GetPages(q Query) ([]int64, int, error) {
	keyList := make([]int64, 0, 100)
	count := 0

	err := mem.db.View(func(tx *memcache.Tx) error {
		b := tx.Bucket([]byte(q.Series))
		if b == nil {
			return fmt.Errorf("Bucket:%v does not exist", q.Series)
		}

		c := b.Cursor()

		first := q.To
		last := q.From
		next := c.Prev
		loopCondition := func(val int64, last int64) bool {
			return val >= last
		}
		if strings.Compare(q.Sort, ASC) == 0 {
			first = q.From
			last = q.To
			next = c.Next
			loopCondition = func(val int64, last int64) bool {
				return val <= last
			}

		}

		// Iterate over the time values
		var k []byte

		for k, _ = c.Seek(timeToByteArr(first)); k != nil && loopCondition(byteArrToTime(k), last); k, _ = next() {
			if count%q.MaxEntries == 0 {
				keyList = append(keyList, byteArrToTime(k))
			}
			count = count + 1
		}
		return nil
	})

	if err != nil {
		return nil, 0, err
	}
	return keyList, count, nil
}

func (mem MemCacheDB) Get(series string) (TimeSeries, error) {
	timeSeries := make([]TimeEntry, 0, 100)
	err := mem.db.View(func(tx *memcache.Tx) error {

		b := tx.Bucket([]byte(series))
		if b == nil {
			return fmt.Errorf("Bucket:%v does not exist", series)
		}

		err := b.ForEach(func(k, v []byte) error {

			record := TimeEntry{byteArrToTime(k), v}
			// TODO: 1. It is an inefficient way of keeping the slices.
			// This has to be addressed during the pagination implementation
			timeSeries = append(timeSeries, record)
			return nil
		})
		if err != nil {
			return err
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	return timeSeries, err
}

func (mem MemCacheDB) GetOnChannel(series string) (<-chan TimeEntry, chan error) {

	resultCh := make(chan TimeEntry, 10)
	errorCh := make(chan error)
	go func() {
		//defer close(resultCh)
		defer close(errorCh)
		err := mem.db.View(func(tx *memcache.Tx) error {

			b := tx.Bucket([]byte(series))
			if b == nil {
				return fmt.Errorf("Bucket:%v does not exist", series)
			}

			err := b.ForEach(func(k, v []byte) error {
				record := TimeEntry{byteArrToTime(k), v}
				resultCh <- record
				return nil
			})
			if err != nil {
				return err
			}
			return err
		})
		// make sure you close the resultchannel before error channel
		close(resultCh)
		if err != nil {
			errorCh <- err
			return
		}
	}()

	return resultCh, errorCh
}

func (mem MemCacheDB) Delete(series string) error {
	return mem.db.Update(func(tx *memcache.Tx) error {
		err := tx.DeleteBucket([]byte(series))
		if err == memcache.ErrBucketNotFound {
			return ErrSeriesNotFound
		}
		return err
	})
}
