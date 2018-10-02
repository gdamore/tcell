# Go logging rethinked

[![Circle CI](https://circleci.com/gh/yanzay/log.svg?style=svg)](https://circleci.com/gh/yanzay/log)
[![Build Status](https://travis-ci.org/yanzay/log.svg?branch=master)](https://travis-ci.org/yanzay/log)
[![Coverage Status](https://coveralls.io/repos/github/yanzay/log/badge.svg?branch=master)](https://coveralls.io/github/yanzay/log?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/yanzay/log)](https://goreportcard.com/report/github.com/yanzay/log)

[![GoDoc](https://godoc.org/github.com/yanzay/log?status.svg)](https://godoc.org/github.com/yanzay/log)

## Usage

```go
log.Println("some log") // unconditional log
log.Trace("trace") // log only with `trace` level
log.Tracef("42: %s", "yep") // each method has it's format alternative
log.Debug("debug") // log only with `debug` level and lower
log.Info("info") // log only with `info` level and lower
log.Warning("warn") // log with `warning` level and lower
log.Error("err") // log with `error` and `fatal` level
log.Fatal("haha") // log and panic("haha")
```

Log adds `--log-level` flag to your program:

```go
// main.go
package main

import (
    "flag"

    "github.com/yanzay/log"
)

func main() {
    flag.Parse()
    log.Info("info")
}
```

```bash
$ go run main.go --help
Usage:
  -log-level string
        Log level: trace|debug|info|warning|error|fatal (default "info")
```

## Advanced Usage

### Log Level

You can set logging level manually by:
```go
log.Level = log.LevelTrace
```

### Log Writers

#### DefaultWriter

DefaultWriter is just a small wrapper for `log` package from stdlib.

#### AsyncWriter

You can use AsyncWriter to write your logs asynchronously, without blocking you main execution. To use AsyncWriter, just switch it in your code like this:

```go
log.Writer = log.NewAsyncWriter()
```

#### Your Own Writer

Also you can use your own log writer:

```go
log.Writer = myWriter // myWriter should implement io.Writer interface
```

