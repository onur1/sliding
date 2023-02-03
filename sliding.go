package sliding

import (
	"math"
	"time"

	"github.com/onur1/ring"
)

// A Counter is a sliding-window based counter which is used for
// counting how many times something have occured within the last time frame.
type Counter struct {
	count chan struct{}
	exit  chan struct{}
	peek  chan chan int
}

// NewCounter returns a new Counter with the given frame duration.
func NewCounter(d time.Duration) *Counter {
	exit := make(chan struct{}, 1)
	count := make(chan struct{})
	peek := make(chan chan int)

	go loop(d, count, exit, peek)

	return &Counter{
		count: count,
		exit:  exit,
		peek:  peek,
	}
}

// Close stops the counter.
func (c *Counter) Close() {
	close(c.exit)
}

// Peek returns the current count.
func (c *Counter) Peek() int {
	res := make(chan int, 1)
	c.peek <- res
	return <-res
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	c.count <- struct{}{}
}

var (
	log2        = math.Log(2)
	minInterval = float64(time.Millisecond.Nanoseconds())
)

func loop(d time.Duration, cnt chan struct{}, exit chan struct{}, peek chan chan int) {
	var (
		ptr  uint = 0
		head int  = 0
	)

	t := d.Nanoseconds()
	l := uint(math.Ceil(math.Log(float64(t)) / log2))

	r := ring.NewRing[int](l + 1)
	ticker := time.NewTicker(
		time.Duration(math.Max(float64(t/int64(l)), minInterval)) * time.Nanosecond,
	)

	var a, b int

LOOP:
	for {
		select {
		case <-cnt:
			head = head + 1
			r.Put(ptr, head)
		case <-ticker.C:
			ptr = ptr + 1
			ptr = r.Put(ptr, head)
		case ch := <-peek:
			a, b = 0, 0
			if r.Get(ptr) != 0 {
				a = r.Get(ptr)
			}
			if r.Get(ptr-l) != 0 {
				b = r.Get(ptr - l)
			}
			ch <- a - b
		case <-exit:
			break LOOP
		}
	}

	ticker.Stop()

	r, ticker = nil, nil

	close(cnt)
	close(peek)
}
