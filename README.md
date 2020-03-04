# go-cacher

cacheしてくれるやつ

## Usage

[_example/main.go](_example/main.go)
```go
package main

import (
	"net/http"
	"time"

	"github.com/nu50218/go-cacher"
)

func work(resp *http.Response) {
	// do something
}

func main() {
	c := cacher.New(5*time.Minute, func() interface{} {
		resp, _ := http.Get("hogehoge")
		return resp
	})

	for i := 0; i < 100; i++ {
		go work(c.Load().(*http.Response))
	}
}
```
