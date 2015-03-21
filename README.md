# goakaibu - Akaibu log archive reader and writer for Go

Supports the full [Akaibu log archive specification version 1](https://github.com/nickbruun/akaibu-format). Documentation is available on [godoc.org](http://godoc.org/github.com/nickbruun/goakaibu).

## Testing

To test against samples, generate the samples with `make samples` in the Akaibu format repository, and place the files in `fixtures/` before running tests. Tests are run with Make:

    make test
