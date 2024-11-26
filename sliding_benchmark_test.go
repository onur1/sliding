package sliding_test

import (
	"sync"
	"testing"
	"time"

	"github.com/onur1/sliding"
	slidingv1 "github.com/onur1/sliding/v1"
)

func BenchmarkSlidingAtomic(b *testing.B) {
	c := sliding.NewCounter(time.Second * 1)

	for i := 0; i < b.N; i++ {
		c.Increment()
	}

	result := c.Peek()
	b.Logf("expected=%d result=%d diff=%d", b.N, result, b.N-int(result))
}

func BenchmarkSlidingChannels(b *testing.B) {
	c := slidingv1.NewCounter(time.Second * 1)

	for i := 0; i < b.N; i++ {
		c.Inc()
	}

	result := c.Peek()
	b.Logf("expected=%d result=%d diff=%d", b.N, result, b.N-int(result))
}

func BenchmarkSlidingConcurrentAtomic(b *testing.B) {
	if b.N == 1 {
		return
	}

	c := sliding.NewCounter(time.Second * 1)

	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/4; i++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()
	result := c.Peek()
	b.Logf("expected=%d result=%d diff=%d", b.N, result, b.N-int(result))
}

func BenchmarkSlidingConcurrentChannels(b *testing.B) {
	if b.N == 1 {
		return
	}

	c := slidingv1.NewCounter(time.Second * 1)

	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/4; i++ {
				c.Inc()
			}
		}()
	}

	wg.Wait()
	result := c.Peek()
	b.Logf("expected=%d result=%d diff=%d", b.N, result, b.N-int(result))
}
