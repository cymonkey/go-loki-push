<h1 align="center">
  go-loki-push
</h1>

<p align="center">
  A lightweight client for pushing logs to Loki written in <a href="https://golang.org/">Golang</a>
</p>

<p align="center">
  <a href="https://cymonkey.mit-license.org/"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
</p>

Inspired by https://pkg.go.dev/github.com/grafana/loki, find that it is a sledgehammer for cracking a nut so I create this package as a simple requests pusher to loki, without any unwanted configs or unused dependencies.

This package just does simple jobs, batching requests and sending them to loki, it provides a grpc/http client pusher to push logs to Loki ingester.

## Installation
```
$ go get github.com/cymonkey/go-loki-push
```

## Usage
```go
import (
    "github.com/cymonkey/go-loki-push/loki"
)

func main() {
    loki.NewClient(NewWithDefaultConfig(&loki.Config{}))
}
```

## License

MIT License, check [LICENSE](./LICENSE).
