
# sliding

**sliding** is a Go library for counting occurrences within a moving time window, allowing efficient tracking of event frequency over a configurable time span. This is useful for scenarios where real-time tracking of event rates is needed, such as API rate limiting, analytics, or monitoring.

## Features

- **Sliding Window Counter**: Counts events within a sliding window of time, providing real-time tracking.
- **Efficient Ring Buffer**: Uses an optimized ring buffer for efficient memory use, with capacities optimized to powers of 2.

## Installation

Install the library by running:

```sh
go get github.com/onur1/sliding
```

## Usage

Here’s an example demonstrating the `sliding` counter:

```go
package main

import (
    "fmt"
    "time"
    "github.com/onur1/sliding"
)

func main() {
    c := sliding.NewCounter(time.Millisecond * 77)

    fmt.Println(c.Peek()) // Output: 0

    c.Inc()
    c.Inc()

    fmt.Println(c.Peek()) // Output: 2

    time.Sleep(time.Millisecond * 55)

    c.Inc()

    fmt.Println(c.Peek()) // Output: 3

    time.Sleep(time.Millisecond * 22)

    fmt.Println(c.Peek()) // Output: 1
}
```

## API Reference

- **`NewCounter(duration time.Duration) *Counter`**: Initializes a new sliding counter with the specified time window.
- **`Inc()`**: Increments the event count by 1.
- **`Peek() int`**: Returns the current count of events within the sliding window.

## License

MIT License. See [LICENSE](LICENSE) for details.
