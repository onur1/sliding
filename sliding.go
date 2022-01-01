package sliding

import (
	"math"
	"time"

	"github.com/onur1/ring"
)

type Counter struct {
	count chan int
	exit  chan struct{}
	peek  chan chan int
}

func (c *Counter) Close() {
	close(c.exit)
}

func (c *Counter) Peek() int {
	res := make(chan int, 1)
	c.peek <- res
	return <-res
}

func (c *Counter) Inc(n int) {
	c.count <- n
}

var (
	log2        = math.Log(2)
	minInterval = float64(time.Millisecond.Nanoseconds())
)

func loop(d time.Duration, cnt chan int, exit chan struct{}, peek chan chan int) {
	var (
		ptr  int = 0
		head int = 0
	)

	t := d.Nanoseconds()
	l := int(math.Ceil(math.Log(float64(t)) / log2))

	list, ticker := ring.NewRing(l+1), time.NewTicker(
		time.Duration(math.Max(float64(t/int64(l)), minInterval))*time.Nanosecond,
	)

LOOP:
	for {
		select {
		case i := <-cnt:
			head = head + i
			list.Put(ptr, head)
		case <-ticker.C:
			ptr = ptr + 1
			ptr = list.Put(ptr, head)
		case ch := <-peek:
			ch <- list.Get(ptr).(int) - (list.Get(ptr - l).(int))
		case <-exit:
			break LOOP
		}
	}

	close(cnt)
	close(peek)
}

func NewCounter(d time.Duration) *Counter {
	exit := make(chan struct{}, 1)
	count := make(chan int)
	peek := make(chan chan int)

	go loop(d, count, exit, peek)

	return &Counter{
		count: count,
		exit:  exit,
		peek:  peek,
	}
}
