package cacher

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_Map(t *testing.T) {
	cnt := int32(0)

	m := NewMap(func(key interface{}) Cacher {
		return New(time.Hour, func() interface{} {
			atomic.AddInt32(&cnt, 1)
			return key
		})
	})

	wg := sync.WaitGroup{}
	fatalChan := make(chan string)

	for i := range [10000]struct{}{} {
		wg.Add(1)
		go func(i int) {
			if val := m.Get(i % 100).Load(); val != i%100 {
				fatalChan <- fmt.Sprintf("unexpected val: %v (want: %v)", val, i%100)
			}
			wg.Done()
		}(i)
	}

	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		waitChan <- struct{}{}
	}()

	select {
	case s := <-fatalChan:
		t.Fatal(s)
	case <-waitChan:
	}

	if cnt != 100 {
		t.Fatal("not cached correctly")
	}
}
