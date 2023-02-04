package sliding

import (
	"math"
	"time"

	"github.com/onur1/ring"
)

// A Counter is a sliding window based counter and used for counting
// how frequently events are happening.
type Counter struct {
	count chan struct{}
	exit  chan struct{}
	peek  chan chan int64
	dur   time.Duration
}

// NewCounter returns a new Counter with the given time interval.
func NewCounter(d time.Duration) *Counter {
	exit := make(chan struct{}, 1)
	count := make(chan struct{})
	peek := make(chan chan int64)

	millis := int(d.Milliseconds())

	// ringsize is always a power of 2 and less than 64
	r := ring.NewRing[int64](int(math.Ceil(math.Log(float64(millis))/math.Log(2))) + 1)

	len := r.Size() - 1

	// time window value needs to be rounded to the nearest millisecond
	// which is divisible by ringsize - 1
	diff := (millis % len)
	if diff <= len/2 {
		millis = millis - diff
	} else {
		millis = millis + (len - diff)
	}

	c := &Counter{
		count: count,
		exit:  exit,
		peek:  peek,
		dur:   time.Duration(millis) * time.Millisecond,
	}

	go c.loop(r, time.Duration(millis/len)*time.Millisecond, len)

	return c
}

// Duration returns the calculated duration of time window.
func (c *Counter) Duration() time.Duration {
	return c.dur
}

// Stop stops a counter.
func (c *Counter) Stop() {
	close(c.exit)
}

// Peek returns the current count.
func (c *Counter) Peek() int64 {
	res := make(chan int64, 1)
	c.peek <- res
	return <-res
}

// Inc increments a counter by 1.
func (c *Counter) Inc() {
	c.count <- struct{}{}
}

func (c *Counter) loop(r *ring.Ring[int64], d time.Duration, len int) {
	var (
		ptr  int   = 0
		head int64 = 0
	)

	ticker := time.NewTicker(d)

LOOP:
	for {
		select {
		case <-c.exit:
			break LOOP
		default:
		}
		select {
		case <-c.count:
			head = head + 1
			r.Put(ptr, head)
		case <-ticker.C:
			ptr = ptr + 1
			ptr = r.Put(ptr, head)
		case ch := <-c.peek:
			ch <- r.Get(ptr) - r.Get(ptr-len)
		}
	}

	ticker.Stop()

	close(c.count)
	close(c.peek)
}
