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

package memory

/*
It implements a low-level key/value database storage engine. It supports fully
serializable transactions, ACID semantics, and lock-free MVCC with multiple
readers and a single writer. The Bhojpur Cache in-memory database storage engine
can be used for projects that want a simple data store without the need to add
large dependencies such as PostgreSQL or MySQL.

The Bhojpur Cache in-memory database storage engine is a single-level, zero-copy,
B+tree data store. It means that the storage engine is optimized for fast read
access and does not require recovery in the event of a system crash. Transactions
which have not finished committing will simply be rolled back in the event of a
system crash.

The design of Bhojpur Cache in-memory database storage is based on Howard Chu's
LMDB database project.

The Bhojpur Cache in-memory database storage currently works on the Windows, macOS,
and Linux operating system.

Basics

There are only a few types in Bhojpur Cache in-memory database storage: DB,
Bucket, Tx, and Cursor. The DB is a collection of Buckets and is represented
by a single file on the disk. A bucket is a collection of unique keys that
are associated with values.

Transactions provide either read-only or read-write access to the database.
Read-only transactions can retrieve key/value pairs and can use Cursors to
iterate over the dataset sequentially. Read-write transactions can create and
delete buckets and can insert and remove keys. Only one read-write transaction
is allowed at a time.

Caveats

The Bhojpur Cache in-memory database storage engine uses a read-only,
memory-mapped data file to ensure that applications cannot corrupt the database.
However, this means that keys and values returned from Bhojpur Cache in-memory
database storage engine cannot be changed. Writing to a read-only byte slice
will cause the Go runtime engine to panic.

Keys and Values retrieved from the in-memory database are only valid for the
life of the transaction. When used outside the transaction, these byte slices
can point to different data or can point to invalid memory which will cause a
panic.

*/
