# Bhojpur Cache - Management Engine

The Bhojpur Cache is a high-performance data caching platform applied within [Bhojpur.NET Platform](https://github.com/bhojpur/platform) for delivering web-scalable applications or services. It could utilize different kinds of data storage engines (e.g., in-memory, file-based) depending on the application's use cases.

## Key Features

- Multi-modal Storage Engines
- ACID transactions
- Web Services APIs

## Getting Started

To install [Bhojpur Cache](https://github.com/bhojpur/cache), use the `go get` command:

```sh
$ go get github.com/bhojpur/cache/...
```

## Introspection Dashboard

To debug the [Bhojpur Cache](https://github.com/bhojpur/cache), you can add an
introspection handler to an HTTP mux and get bettervisibility into in-memory
database storage engine's behaviour. For example

```sh
$ go build -o bin/cachedbg ./internal/main.go
```

```go
http.Handle("/introspect", http.StripPrefix("/introspect", debugger.NewHandler(mydb)))
```

then, run the `bin/cachedbg` binary by passing in the path to your database:

```sh
$ bin/cachedbg ./internal/path/to/my.db
```

After pointing your web browser to `http://localhost:3000`, you should see something like this 
![Introspection Dashboard](/internal/debugger.png "Bhojpur Cache - In-Memory Database")

It allows you to introspect [Bhojpur Cache](https://github.com/bhojpur/cache)
database in a web browser. The `bin/cachedbg` tool gives you access to
low-level page information and b-tree structures so you can better understand
how [Bhojpur Cache](https://github.com/bhojpur/cache) is laying out your data.

## HTTP Integration

You can also use boltd as an `http.Handler` in your own application. To use it,
simply add the handler to your muxer:

To generate a custom web template, you need the following tool

```sh
$ go get github.com/benbjohnson/ego
```

## Distributed Applications

- [Bhojpur CMS](https://github.com/bhojpur/cms) content management system
- [Bhojpur Graph](https://github.com/bhojpur/graph) graph database system
- [Bhojpur Keyed](https://github.com/bhojpur/keyed) state consensus system
- [Bhojpur SQL](https://github.com/bhojpur/sql) relational database system  
- [Bhojpur UFS](https://github.com/bhojpur/ufs) object storage system
- [Bhojpur Web](https://github.com/bhojpur/web) server and RESTful APIs engine
