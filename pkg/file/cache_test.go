package file_test

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
	"math"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bhojpur/cache/pkg/file"
	"github.com/stretchr/testify/assert"
)

func Example_Cache() {
	cache := file.New("./cache.data")
	defer os.Remove("./cache.data")

	v, err := cache.Get("not-exist")
	fmt.Println(v.Valid(), err)

	fmt.Println(cache.Set("k", "v", time.Minute))

	v, err = cache.Get("k")
	fmt.Println(v.Valid(), v.String(), err)

	ttl, err := cache.TTL("k")
	fmt.Println(int(math.Ceil(ttl.Seconds())), err)

	time.Sleep(time.Second)

	ttl, err = cache.TTL("k")
	fmt.Println(int(math.Ceil(ttl.Seconds())), err)

	fmt.Println(cache.Del("k"))

	v, err = cache.Get("k")
	fmt.Println(v.Valid(), err)

	// output:
	// false <nil>
	// <nil>
	// true v <nil>
	// 60 <nil>
	// 59 <nil>
	// <nil>
	// false <nil>
}

func TestNew(t *testing.T) {
	as := assert.New(t)
	defer os.Remove("./test")

	os.Remove("./test")
	c := file.New("./test")

	t.Run("not found", func(t *testing.T) {
		v, err := c.Get("k1")
		as.Equal(nil, err)
		as.False(v.Valid())
		as.Empty(v.String())
	})

	t.Run("exist get set", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		v, err := c.Get("k")
		as.Nil(err)
		as.Equal("v", v.String())
	})

	t.Run("expired", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		time.Sleep(time.Second)

		v, err := c.Get("k")
		as.Equal(nil, err)
		as.False(v.Valid())
		as.Empty(v.String())
	})

	t.Run("ttl", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		ttl, err := c.TTL("k")
		as.Nil(err)
		as.True(ttl <= time.Second && ttl >= time.Second-100*time.Millisecond, ttl)
	})

	t.Run("expire", func(t *testing.T) {
		as.Nil(c.Set("k", "v", time.Second))

		ttl, err := c.TTL("k")
		as.Nil(err)
		as.True(ttl <= time.Second && ttl >= time.Second-100*time.Millisecond)

		as.Nil(c.Expire("k", time.Minute))

		ttl, err = c.TTL("k")
		as.Nil(err)
		as.True(ttl <= time.Minute && ttl >= time.Minute-100*time.Millisecond)
	})

	t.Run("range", func(t *testing.T) {
		os.Remove("./test")
		c = file.New("./test")

		for i := 0; i < 1000; i++ {
			j := strconv.Itoa(i)
			as.Nil(c.Set(j, j, time.Minute), i)
		}

		kvs, err := c.Range()
		as.Nil(err)
		for _, v := range kvs {
			as.Equal(v.Key, v.Val)
		}
		as.Len(kvs, 1000)
	})
}
