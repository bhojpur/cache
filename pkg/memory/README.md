# Bhojpur Cache - In-Memory Storage Engine

It is a key/value database `storage engine` inspired by [Howard Chu's][hyc_symas]
[LMDB project][lmdb]. The goal of the project is to provide a simple, fast, and
reliable in-memory database storage engine for such projects that do not
require full-fledged database server functionality, such as: PostgreSQL or MySQL.

Since the [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database
storage engine is meant to be used as a low-level piece of functionality, thence
simplicity is the key. The Database APIs will be small and only focus on getting
values and setting values.

## Getting Started

### Installing

To start using [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database
storage engine, install Go and run `go get`:

```sh
$ go get github.com/bhojpur/cache/pkg/memory...
```

It will retrieve the `in-memory storage` database engine library and install the
[Bhojpur Cache](https://github.com/bhojpur/cache) command line utility into
your `$GOBIN` path.


### Opening an In-Memory Database

The top-level object in a [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory
database storage engine is a `DB`. It is represented as a single `file` on your data
storage volume and represents a consistent `snapshot` of your in-memory data.

To open your `in-memory database`, simply use the `memory.Open()` function:

```go
package main

import (
	"log"

	memory "github.com/bhojpur/cache/pkg/memory"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created, if the file doesn't exist.
	db, err := memory.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	...
}
```

Please note that [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database
storage engine obtains a **file lock** on the data file so multiple processes cannot
open the same database at the same time. Opening an already open [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database file will cause it to hang until the other processes closes
it. To prevent an indefinite wait time, you can pass a `timeout` option to the `Open()`
function:

```go
db, err := memory.Open("my.db", 0600, &memory.Options{Timeout: 1 * time.Second})
```


### In-Memory Transactions

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine allows only `one read-write` transaction at a time, but allows as `many read-only`
transactions as you want at a time. Each transaction has a consistent view of
the data as it existed when the transaction started.

Individual transactions and all `objects` created from them (e.g. buckets, keys)
are not `thread safe`. To work with data in multiple `goroutines` you must start
a transaction for each one or use `locking` to ensure only one goroutine accesses
a transaction at a time. Creating a transaction from the `DB` is `thread safe`.

The `read-only` transactions and `read-write` transactions should not depend on
one another and generally shouldn't be opened simultaneously in the same goroutine.
It can cause a deadlock as the `read-write` transaction needs to periodically
re-map the data file, but it cannot do so while a `read-only` transaction is open.


#### Read-write Transactions

To start a `read-write` transaction, you can use the `DB.Update()` function:

```go
err := db.Update(func(tx *memory.Tx) error {
	...
	return nil
})
```

Inside the closure, you have a consistent view of the database. You `commit` the
transaction by returning `nil` at the end. You can also `rollback` the transaction
at any point by returning an error. All database operations are allowed inside
a `read-write` transaction.

Always check the return `error` as it will report any disk failures that can cause
your transaction to remain incomplete. If you return an `error` within your closure
it will be passed through.

#### Read-only Transactions

To start a `read-only` transaction, you can use the `DB.View()` function:

```go
err := db.View(func(tx *memory.Tx) error {
	...
	return nil
})
```

You also get a consistent view of the database within this closure. However,
no mutating operations are allowed within a `read-only` transaction. You can
only retrieve the `buckets`, retrieve values, and copy the database within a
`read-only` transaction.


#### Batch read-write Transactions

Each `DB.Update()` operation waits for the storage disk volumes to commit the
`writes`. This overhead can be minimized by combining multiple updates with
the `DB.Batch()` function:

```go
err := db.Batch(func(tx *memory.Tx) error {
	...
	return nil
})
```

The concurrent `Batch` calls are opportunistically combined into larger
transactions. A `Batch` is only useful when there are multiple goroutines
calling it.

The trade-off is that `Batch` can call the given function multiple times,
if parts of the transaction fail. The function must be idempotent and side
effects must take effect only after a successful return from `DB.Batch()`.

For example: do not display messages from inside the function, instead
set variables in the enclosing scope:

```go
var id uint64
err := db.Batch(func(tx *memory.Tx) error {
	// Find last key in bucket, decode as bigendian uint64, increment
	// by one, encode back to []byte, and add new key.
	...
	id = newValue
	return nil
})
if err != nil {
	return ...
}
fmt.Println("Allocated ID %d", id)
```


#### Managing transactions manually

The `DB.View()` and `DB.Update()` functions are wrappers around the `DB.Begin()`
function. These helper functions will start the transaction, execute a function,
then safely close your transaction if an error is returned. It is a recommended
way to use [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database
transactions.

However, sometimes you may want to manually `start` and `end` your transactions.
You can use the `DB.Begin()` function directly, but **please** be sure to close
the transaction.

```go
// Start a writable transaction.
tx, err := db.Begin(true)
if err != nil {
    return err
}
defer tx.Rollback()

// Use the transaction...
_, err := tx.CreateBucket([]byte("MyBucket"))
if err != nil {
    return err
}

// Commit the transaction and check for error.
if err := tx.Commit(); err != nil {
    return err
}
```

The first argument to `DB.Begin()` is a `boolean` stating, if the transaction
should be `writable`.


### Using Buckets

The `Bucket` are collections of key/value pairs within the database. All keys
in a bucket must be unique. You can create a `Bucket` using the `DB.CreateBucket()`
function:

```go
db.Update(func(tx *memory.Tx) error {
	b, err := tx.CreateBucket([]byte("MyBucket"))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
})
```

You can also create a `Bucket` only if it doesn't exist by using the
`Tx.CreateBucketIfNotExists()` function. It's a common pattern to call this
function for all your top-level buckets after you open your database so that
you can guarantee they exist for future transactions.

To delete a `Bucket`, simply call the `Tx.DeleteBucket()` function.


### Using key/value Pairs

To save a key/value pair to a `Bucket`, use the `Bucket.Put()` function:

```go
db.Update(func(tx *memory.Tx) error {
	b := tx.Bucket([]byte("MyBucket"))
	err := b.Put([]byte("answer"), []byte("42"))
	return err
})
```

It will set the value of the `"answer"` key to `"42"` in the `MyBucket`
bucket. To retrieve this value, we can use the `Bucket.Get()` function:

```go
db.View(func(tx *memory.Tx) error {
	b := tx.Bucket([]byte("MyBucket"))
	v := b.Get([]byte("answer"))
	fmt.Printf("The answer is: %s\n", v)
	return nil
})
```

The `Get()` function does not return any error, because its operation is
guaranteed to work (unless there is some kind of system failure). If the `key`
exists, then it will return its byte slice value. If it doesn't exist, then it
will return `nil`. It is important to note that you can have a zero-length
value set to a `key` which is different than the key not existing.

Use the `Bucket.Delete()` function to delete a key from the `Bucket`.

Please note that values returned from the `Get()` are only valid, while the
transaction is open. If you need to use a value outside of the transaction
then you must use `copy()` to copy it to another byte slice.


### Auto-incrementing integer for the bucket

By using the `NextSequence()` function, you can let [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database storage engine determine a `sequence`, which can be used as the
unique identifier for your key/value pairs. See the example below.

```go
// CreateUser saves u to the In-Memory database. The new user ID is set on u once the data is persisted.
func (s *Store) CreateUser(u *User) error {
    return s.db.Update(func(tx *memory.Tx) error {
        // Retrieve the users Bucket.
        // This should be created when the In-Memory database is first opened.
        b := tx.Bucket([]byte("users"))

        // Generate ID for the user.
        // This returns an error only if the Tx is closed or not writeable.
        // That can't happen in an Update() call so I ignore the error check.
        id, _ := b.NextSequence()
        u.ID = int(id)

        // Marshal user data into bytes.
        buf, err := json.Marshal(u)
        if err != nil {
            return err
        }

        // Persist bytes to users Bucket.
        return b.Put(itob(u.ID), buf)
    })
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

type User struct {
    ID int
    ...
}
```

### Iterating over Keys

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine stores its `keys` in byte-sorted order within a `Bucket`. It makes sequential
iteration over these `keys` extremely fast. To iterate over the `keys`, we'll use a
`Cursor`:

```go
db.View(func(tx *memory.Tx) error {
	// Assuming that Bucket exists and has keys
	b := tx.Bucket([]byte("MyBucket"))

	c := b.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("key=%s, value=%s\n", k, v)
	}

	return nil
})
```

The `Cursor` allows you to move to a specific point in the list of `keys` and
move `forward` or `backward` through the keys one at a time.

The following functions are available on the `Cursor` object:

```
First()  Move to the first key.
Last()   Move to the last key.
Seek()   Move to a specific key.
Next()   Move to the next key.
Prev()   Move to the previous key.
```

Each of those functions has a return signature of `(key []byte, value []byte)`.
When you have iterated to the end of the `Cursor`, then `Next()` will return a
`nil` key. You must seek to a position using `First()`, `Last()`, or `Seek()`
before calling `Next()` or `Prev()`. If you do not seek to a position, then
these functions will return a `nil` key.

During iteration, if the `key` is non-`nil` but the value is `nil`, that means
the `key` refers to a `Bucket` rather than a value.  Use `Bucket.Bucket()` to
access the sub-bucket.

#### Prefix Scans

To iterate over a `key` prefix, you can combine `Seek()` and `bytes.HasPrefix()`:

```go
db.View(func(tx *memory.Tx) error {
	// Assuming that Bucket exists and has keys
	c := tx.Bucket([]byte("MyBucket")).Cursor()

	prefix := []byte("1234")
	for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
		fmt.Printf("key=%s, value=%s\n", k, v)
	}

	return nil
})
```

#### Range Scans

Another common use case is scanning over a `range` such as, a `time range`. If
you use a sortable time encoding, such as `RFC3339`, then you can query a
specific `date range` like this:

```go
db.View(func(tx *memory.Tx) error {
	// Assume our events bucket exists and has RFC3339 encoded time keys.
	c := tx.Bucket([]byte("Events")).Cursor()

	// Our time range spans the 2010's decade.
	min := []byte("2010-01-01T00:00:00Z")
	max := []byte("2020-01-01T00:00:00Z")

	// Iterate over the 2010's.
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		fmt.Printf("%s: %s\n", k, v)
	}

	return nil
})
```

Note that, while `RFC3339` is sortable, the Go implementation of `RFC3339Nano`
does not use a fixed number of digits after the decimal point and is therefore
not sortable.

#### ForEach()

You can also use the function `ForEach()`, if you know you'll be iterating over
all the `keys` in a Bucket:

```go
db.View(func(tx *memory.Tx) error {
	// Assume that Bucket exists and has keys
	b := tx.Bucket([]byte("MyBucket"))

	b.ForEach(func(k, v []byte) error {
		fmt.Printf("key=%s, value=%s\n", k, v)
		return nil
	})
	return nil
})
```

Please note that `keys` and `values` in a `ForEach()` call are only valid, while
the transaction is `open`. If you need to use a `key` or `value` outside of the
transaction, you must use `copy()` to copy it to another byte slice.

### Nested Buckets

You can also store a `Bucket` in a key to create nested buckets. The API is the
same as the Bucket Management API on the `DB` object:

```go
func (*Bucket) CreateBucket(key []byte) (*Bucket, error)
func (*Bucket) CreateBucketIfNotExists(key []byte) (*Bucket, error)
func (*Bucket) DeleteBucket(key []byte) error
```

For example, you had a `multi-tenant` software application, where the root-level
bucket was the `Account` bucket. Inside of this bucket, there was a sequence of
`accounts`, which themselves are buckets. And, inside the sequence bucket, you
could have many more buckets pertaining to the `Account` itself (e.g., Users,
Notes, etc) isolating the information into logical groupings.

```go

// createUser creates a new user in the given account.
func createUser(accountID int, u *User) error {
    // Start the In-Memory database transaction.
    tx, err := db.Begin(true)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Retrieve the root Bucket for the account.
    // Assume this has already been created when the account was set up.
    root := tx.Bucket([]byte(strconv.FormatUint(accountID, 10)))

    // Setup the users Bucket.
    bkt, err := root.CreateBucketIfNotExists([]byte("USERS"))
    if err != nil {
        return err
    }

    // Generate an ID for the new User.
    userID, err := bkt.NextSequence()
    if err != nil {
        return err
    }
    u.ID = userID

    // Marshal and save the encoded User.
    if buf, err := json.Marshal(u); err != nil {
        return err
    } else if err := bkt.Put([]byte(strconv.FormatUint(u.ID, 10)), buf); err != nil {
        return err
    }

    // Commit the In-Memory database transaction.
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

```

### In-Memory Database Backups

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine stores a single `file` so it's easy to backup. You can use `Tx.WriteTo()`
function to write a consistent view of the in-memory database to a writer. If you
call this from a `read-only` transaction, it will perform a `hot backup` and not
block your other database reads and writes.

By default, it will use a regular file handle which will utilize the operating
system's page cache.

A common use case is to `Backup-over-HTTP` so that you could use tools like `cURL`
to do the in-memory database backups:

```go
func BackupHandleFunc(w http.ResponseWriter, req *http.Request) {
	err := db.View(func(tx *memory.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

Then, you can `backup` the data using this command:

```sh
$ curl http://localhost/backup > my.db
```

Or, you can open a web-browser and point to `http://localhost/backup`. It will
download the in-memory data `snapshot` automatically.

If you want to backup to another file you can use the `Tx.CopyFile()` helper
function.


### Statistics

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database keeps
a running count of many of the internal operations as it performs so that you
can better understand what's going on. By grabbing an in-memory data `snapshot`
of these stats at two points in time, we can analyze what operations were
performed during that `time range`.

For example, we could start a goroutine to log the `stats` every 10 seconds:

```go
go func() {
	// Grab the initial stats.
	prev := memdb.Stats()

	for {
		// Wait for 10s.
		time.Sleep(10 * time.Second)

		// Grab the current stats and diff them.
		stats := memdb.Stats()
		diff := stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		json.NewEncoder(os.Stderr).Encode(diff)

		// Save stats for the next loop.
		prev = stats
	}
}()
```

It's also useful to `pipe` these stats to a service, such as: `statsd`, for
monitoring or to provide an HTTP endpoint that would perform a fixed-length
sample.

### Read-only Mode

Sometimes it is useful to create a shared, `read-only` [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database. To do this, set the `Options.ReadOnly` flag when opening
your in-memory database. The `read-only` mode uses a shared lock to allow
multiple processes to read from the in-memory database, but it will block
any processes from opening the database file in `read-write` mode.

```go
db, err := memory.Open("my.db", 0666, &memory.Options{ReadOnly: true})
if err != nil {
	log.Fatal(err)
}
```

### Mobile Platform usage (e.g., Android / iOS)

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine is able to run on mobile devices by leveraging the binding feature of the
[GoMobile](https://github.com/golang/mobile) tool. Create a `struct` that will
contain your [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database
storage logic and a reference to a `*memory.DB` by initializing constructor that
takes in a file path where the database file will be stored. Neither the `Android`
nor `iOS` require extra permissions or cleanup from using this method.

```go
func NewCacheDB(filepath string) *CacheDB {
	db, err := memory.Open(filepath+"/demo.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &CacheDB{db}
}

type CacheDB struct {
	db *memory.DB
	...
}

func (b *CacheDB) Path() string {
	return b.db.Path()
}

func (b *CacheDB) Close() {
	b.db.Close()
}
```

The database logic should be defined as `methods` on this wrapper `struct`.

To initialize this `struct` from the native language (both the mobile platforms
now sync their local storage to the Cloud. These snippets disable that
functionality for the database file):

#### Android Platform

```java
String path;
if (android.os.Build.VERSION.SDK_INT >=android.os.Build.VERSION_CODES.LOLLIPOP){
    path = getNoBackupFilesDir().getAbsolutePath();
} else{
    path = getFilesDir().getAbsolutePath();
}
Cachemobiledemo.CacheDB cacheDB = Cachemobiledemo.NewCacheDB(path)
```

#### iOS Platform

```objc
- (void)demo {
    NSString* path = [NSSearchPathForDirectoriesInDomains(NSLibraryDirectory,
                                                          NSUserDomainMask,
                                                          YES) objectAtIndex:0];
	GoCachemobiledemoCacheDB * demo = GoCachemobiledemoNewCacheDB(path);
	[self addSkipBackupAttributeToItemAtPath:demo.path];
	//Some DB Logic would go here
	[demo close];
}

- (BOOL)addSkipBackupAttributeToItemAtPath:(NSString *) filePathString
{
    NSURL* URL= [NSURL fileURLWithPath: filePathString];
    assert([[NSFileManager defaultManager] fileExistsAtPath: [URL path]]);

    NSError *error = nil;
    BOOL success = [URL setResourceValue: [NSNumber numberWithBool: YES]
                                  forKey: NSURLIsExcludedFromBackupKey error: &error];
    if(!success){
        NSLog(@"Error excluding %@ from backup %@", [URL lastPathComponent], error);
    }
    return success;
}

```

## Comparing with other Database Systems

### PostgreSQL, MySQL, & other relational databases

Relational databases structure data into rows and are only accessible through
the use of SQL. This approach provides flexibility in how you store and query
your data, but also incurs overhead in parsing and planning SQL statements. The
[Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database engine
accesses all data by a byte slice key. This [Bhojpur Cache](https://github.com/bhojpur/cache)
memory database fast to read and write data by key, but provides no built-in
support for joining values together.

Most relational databases (with the exception of `SQLite`) are standalone servers
that run separately from your application. This gives your systems flexibility
to connect multiple application servers to a single database server, but also
adds overhead in serializing and transporting data over the network. The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database engine runs as a library included in your
application, so all data access has to go through your application's process.
This brings data closer to your application, but limits multi-process access
to the data.

### LevelDB, RocksDB

The `LevelDB` and its derivatives (e.g., RocksDB, HyperLevelDB) are similar to
the [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine in that they are libraries bundled into the application, however, their
underlying structure is a log-structured merge-tree (LSM tree). An `LSM` tree
optimizes random writes by using a `write ahead` log and multi-tiered, sorted
files, called `SSTables`. The [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database storage engine uses a `B+tree` internally and only a single
file. Both the approaches have some trade-offs.

If you require a high random write throughput (>10,000 w/sec) or you need to use
spinning disk drives, then `LevelDB` could be a good choice. If your application
is `read-heavy` or does a lot of range scans then [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database storage engine could be a good choice.

Another important consideration is that `LevelDB` does not have transactions.
It supports batch writing of key/values pairs and it supports read snapshots
but it will not give you the ability to do a `compare-and-swap` operation safely.
The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine supports fully serializable ACID transactions.

### LMDB

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory databse storage
engine was originally a port of `LMDB` so it is architecturally similar. Both use
a `B+tree`, have ACID semantics with fully serializable transactions, and support
lock-free MVCC using a `single writer` and `multiple readers`.

The two projects have somewhat diverged. LMDB heavily focuses on raw performance
while [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine has focused on simplicity and ease of use. For example, LMDB allows several
unsafe actions, such as: `direct writes` for the sake of performance. The
[Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage engine
opts to disallow actions which can leave the database in a corrupted state. The
only exception to this in [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database storage engine is `DB.NoSync`.

There are also a few differences in API. LMDB requires a maximum `mmap` size when
opening an `mdb_env` whereas [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database storage engine will handle incremental `mmap` resizing
automatically. LMDB overloads the `getter` and `setter` functions with
multiple flags, whereas [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database splits these specialized cases into their own functions.

## Caveats & Limitations

It's important to pick the right tool for the job and [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage engine is no exception. Here are a few things to
note when evaluating and using it:

* [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
  engine is good for read intensive workloads. Sequential write performance is
  also fast but random writes can be slow. You can use `DB.Batch()` or add a
  write-ahead log to help mitigate this issue.

* [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
  engine uses a `B+tree` internally, so there can be a lot of random page access.
  The **solid-state drive** (SSD) provide a significant performance boost over
  spinning disk drives.

* Try to avoid `long running` read transactions. [Bhojpur Cache](https://github.com/bhojpur/cache)
  in-memory database storage engine uses `copy-on-write` so old pages cannot be
  reclaimed while an old transaction is using them.

* Byte slices returned from [Bhojpur Cache](https://github.com/bhojpur/cache)
  in-memory database storage engine are only valid during a transaction. Once
  the transaction has been committed or rolled back then the memory they point
  to can be reused by a new page or can be unmapped from virtual memory and
  you'll see an `unexpected fault address` panic when accessing it.

* [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
  engine uses an exclusive write lock on the database file so it cannot be
  shared by multiple processes.

* Be careful while using `Bucket.FillPercent`. Setting a high fill percent for
  the `Buckets` that have random inserts will cause your database to have very
  poor page utilization.

* In general, use **larger** buckets. Smaller buckets causes poor memory page
  utilization once they become larger than the page size (typically 4KB).

* Bulk loading a lot of random writes into a new `Bucket` could be slow as the
  page will not split until the transaction is committed. Randomly inserting
  more than 100,000 key/value pairs into a single new `Bucket` in a single
  transaction is not advised.

* [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
  engine uses a `memory-mapped` file, so the underlying operating system handles
  the caching of the data. Typically, the OS will cache as much of the file as
  it can in the memory and will release the memory as needed to other processes.
  This means that [Bhojpur Cache](https://github.com/bhojpur/cache) storage engine
  can show very high memory usage when working with large databases. However, this
  is expected and the OS will release memory as needed. [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory storage engine can handle databases much larger than the available
  physical RAM, provided its `memory-map` fits in the process virtual address
  space. It may be problematic on 32-bits systems.

* The data structures in the [Bhojpur Cache](https://github.com/bhojpur/cache)
  in-memory database are memory mapped so the data file will be `endian` specific.
  This means that you cannot copy a [Bhojpur Cache](https://github.com/bhojpur/cache)
  database file from a little endian machine to a big endian machine and have it work.
  For most users this is not a concern since most modern CPUs are little endian.

* Because of the way `pages` are laid out on disk, the [Bhojpur Cache](https://github.com/bhojpur/cache)
  in-memory database storage engine cannot truncate data files and return free pages
  back to the disk. Instead, [Bhojpur Cache](https://github.com/bhojpur/cache)
  in-memory database storage engine maintains a `free list` of unused pages within
  its data file. These free pages can be reused by later transactions. This works
  well for many use cases as the databases generally tend to grow. However, it's
  important to note that deleting large chunks of data will not allow you to
  reclaim that space on disk.

## Reading the Source Code

The [Bhojpur Cache](https://github.com/bhojpur/cache) in-memory database storage
engine is a relatively small code base (<3KLOC) for an embedded, serializable,
transactional key/value database so it can be a good starting point for people
interested in how databases work.

The best places to start are the main entry points into [Bhojpur Cache](https://github.com/bhojpur/cache)
in-memory database storage engine:

- `Open()` - Initializes the reference to the database. It's responsible for
  creating the database if it doesn't exist, obtaining an exclusive lock on the
  file, reading the meta pages, and memory-mapping the file.

- `DB.Begin()` - Starts a read-only or read-write transaction depending on the
  value of the `writable` argument. This requires briefly obtaining the **meta**
  lock to keep track of open transactions. Only one read-write transaction can
  exist at a time so the **rwlock** is acquired during the life of a read-write
  transaction.

- `Bucket.Put()` - Writes a key/value pair into a `Bucket`. After validating the
  arguments, a cursor is used to traverse the B+tree to the page and position
  where they key & value will be written. Once the position is found, the bucket
  materializes the underlying page and the page's parent pages into memory as
  "nodes". These nodes are where mutations occur during read-write transactions.
  These changes get flushed to disk during commit.

- `Bucket.Get()` - Retrieves a key/value pair from a `Bucket`. This uses a cursor
  to move to the page & position of a key/value pair. During a `read-only`
  transaction, the key and value data is returned as a direct reference to the
  underlying mmap file so there's no allocation overhead. For the `read-write`
  transactions, this data may reference the mmap file or one of the in-memory
  node values.

- `Cursor` - This object is simply for traversing the B+tree of on-disk pages
  or in-memory nodes. It can seek to a specific key, move to the first or last
  value, or it can move forward or backward. The cursor handles the movement up
  and down the B+tree transparently to the end user.

- `Tx.Commit()` - Converts the in-memory dirty nodes and the list of free pages
  into pages to be written to disk. Writing to disk then occurs in two phases.
  First, the dirty pages are written to disk and an `fsync()` occurs. Secondly,
  a new meta page with an incremented transaction ID is written and another
  `fsync()` occurs. This `two phase write` ensures that partially written data
  pages are ignored in the event of a crash since the meta page pointing to them
  is never written. Partially written meta pages are invalidated, because they
  are written with a checksum.
