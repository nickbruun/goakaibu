# goakaibu - Akaibu log archive reader and writer for Go

[![Build status](https://travis-ci.org/nickbruun/goakaibu.svg?branch=master)](https://travis-ci.org/nickbruun/goakaibu) [![GoDoc](https://godoc.org/github.com/nickbruun/goakaibu?status.svg)](https://godoc.org/github.com/nickbruun/goakaibu)

Supports the full [Akaibu log archive specification version 1](https://github.com/nickbruun/akaibu-format). Documentation is available on [GoDoc](http://godoc.org/github.com/nickbruun/goakaibu).

## Testing

To test against samples, generate the samples with `make samples` in the Akaibu format repository, and place the files in `fixtures/` before running tests. Tests are run with Make:

    make test
