package cacher

import (
	"context"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func Test_basicCacher_Cache(t *testing.T) {
	cnt := int32(0)

	c := New(100*time.Millisecond, func() interface{} {
		return atomic.AddInt32(&cnt, 1)
	})

	attack := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				runtime.Goexit()
			default:
				c.Load()
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := 0; i < 100; i++ {
		go attack(ctx)
	}
	<-ctx.Done()

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
