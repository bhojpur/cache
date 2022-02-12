package main

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
	"flag"
	"fmt"
	"log"
	"net/http"

	debugger "github.com/bhojpur/cache/pkg/debugger"
	memcache "github.com/bhojpur/cache/pkg/memory"
)

func main() {
	log.SetFlags(0)
	var (
		addr = flag.String("addr", ":3000", "bind address")
	)
	flag.Parse()

	// Validate parameters.
	var path = flag.Arg(0)
	if path == "" {
		log.Fatal("path to your in-memory database file is required")
	}

	// Open the database.
	db, err := memcache.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Enable logging.
	log.SetFlags(log.LstdFlags)

	// Setup the HTTP handlers.
	http.Handle("/", debugger.NewHandler(db))

	// Start the HTTP server.
	go func() { log.Fatal(http.ListenAndServe(*addr, nil)) }()

	fmt.Printf("Bhojpur Cache server listening on http://localhost%s\n", *addr)
	select {}
}
