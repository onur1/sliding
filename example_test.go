package sliding_test

import (
	"fmt"
	"time"

	"github.com/onur1/sliding"
)

func ExampleCounter() {
	c := sliding.NewCounter(time.Millisecond * 77)

	fmt.Println(c.Peek()) // 0

	c.Inc()
	c.Inc()

	fmt.Println(c.Peek()) // 2

	time.Sleep(time.Millisecond * 55)

	c.Inc()

	fmt.Println(c.Peek()) // 3

	time.Sleep(time.Millisecond * 22)

	fmt.Println(c.Peek()) // 1

	// Output:
	// 0
	// 2
	// 3
	// 1
}
