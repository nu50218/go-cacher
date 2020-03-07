package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nu50218/go-cacher"
)

func workWithBody(body []byte) error {
	// do something
	return nil
}

func work(c cacher.Cacher) error {
	val := c.Load()
	if err, ok := val.(error); ok {
		return err
	}
	body := val.([]byte)
	return workWithBody(body)
}

func main() {
	c := cacher.New(5*time.Minute, func() interface{} {
		resp, err := http.Get("hogehoge")
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return body
	})

	for i := 0; i < 100; i++ {
		go work(c)
	}
}
