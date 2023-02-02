package sliding

import (
	"math"
	"time"

	"github.com/onur1/ring"
)

type Counter struct {
	count chan uint
	peek  chan chan uint
	exit  chan struct{}
}

func (c *Counter) Close() {
	close(c.exit)
}

func (c *Counter) Peek() uint {
	res := make(chan uint, 1)
	c.peek <- res
	return <-res
}

func (c *Counter) Inc(n uint) {
	c.count <- n
}

func loop(d time.Duration, cnt chan uint, exit chan struct{}, peek chan chan uint) {
	var (
		ptr  uint = 0
		head uint = 0
	)

	var a, b uint

	millis := uint(d.Milliseconds())
	len := uint(math.Ceil(math.Log(float64(millis)) / math.Log(2)))

	r := ring.NewRing[uint](len + 1)

	ticker := time.NewTicker(
		time.Duration(math.Max(float64((millis/len)), 250)) * time.Millisecond,
	)

LOOP:
	for {
		select {
		case i := <-cnt:
			head = head + i
			r.Put(ptr, head)
		case <-ticker.C:
			ptr = ptr + 1
			ptr = r.Put(ptr, head)
		case ch := <-peek:
			a, b = 0, 0
			if r.Get(ptr) != 0 {
				a = r.Get(ptr)
			}
			if r.Get(ptr-len) != 0 {
				b = r.Get(ptr - len)
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

func NewCounter(d time.Duration) *Counter {
	exit := make(chan struct{}, 1)
	count := make(chan uint)
	peek := make(chan chan uint)

	go loop(d, count, exit, peek)

	return &Counter{
		count: count,
		exit:  exit,
		peek:  peek,
	}
}
