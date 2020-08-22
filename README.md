# anda-server

## Prerequisites

Required:
- Go v1.14+
- PostgreSQL v12.1+

Optional:
- golangci-lint v1.30.0

These are the versions I am using, so earlier versions may be fine.

## Install

### Get source code

Using `go get`:

```bash
$ go get -u github.com/aewens/anda-server
$ cd $GOPATH/src/github.com/aewens/anda-server
```

Manually by hand:

```bash
$ cd $GOPATH/src
$ mkdir -p github.com/aewens
$ cd github.com/aewens
$ git clone github.com/aewens/anda-server
$ cd anda-server
```

### Configure

Fill out the file with database connection info for PostgreSQL.

```bash
$ cp etc/config.orig.json etc/config.json
$ vim etc/config.json # Fill out database credentials
```

### Setup

Using build and run scripts:

```bash
$ scripts/00_build_and_run.sh
```

Manually by hand:

```bash
$ go build -o bin/anda.o && bin/anda.o -config etc/config.json
```

The build and run script adds a few niceties like removing the previous instance and logging to a file, but is not required.
