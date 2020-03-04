package cacher

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_basicCacher_Cache(t *testing.T) {
	cnt := int32(0)

	c := New(100*time.Millisecond, func() interface{} {
		return atomic.AddInt32(&cnt, 1)
	})

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
		L:
			for {
				select {
				case <-ctx.Done():
					break L
				default:
					c.Load()
				}
			}
			wg.Done()
		}(ctx)
	}
	wg.Wait()

	if cnt < 45 || 55 < cnt {
		t.Fatalf("get was called %v times (expected in [45,55])", cnt)
	}
}

func Test_basicCacher_Load(t *testing.T) {
	c := New(100*time.Second, func() interface{} {
		return 1
	})
	if !reflect.DeepEqual(1, c.Load()) {
		t.Fail()
	}
}
