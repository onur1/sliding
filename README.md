# sliding

**sliding** is a lock-free sliding window counter implementation using atomic CAS operations.

## Installation

Install the library by running:

```sh
go get github.com/tetsuo/sliding
```

## API Reference

### `NewCounter`

```go
func NewCounter(windowSize time.Duration) *Counter
```

Initializes a new sliding counter with the specified time window.

* `windowSize` defines the total sliding window duration.
* The number of slots is calculated automatically for optimal performance.

### `Increment`

```go
func (c *Counter) Increment()
```

Increments the count in the current time slot.

### `Peek`

```go
func (c *Counter) Peek() uint64
```

Returns the total count across the sliding window.

### `FrameDuration`

```go
func (c *Counter) FrameDuration() time.Duration
```

Returns the duration of each slot in the sliding window.

## Benchmark

The atomic-based implementation significantly outperforms the channel-based approach in both single-threaded and concurrent scenarios 👍

```
goos: darwin
goarch: amd64
pkg: github.com/tetsuo/sliding
cpu: Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz
BenchmarkSlidingAtomic
BenchmarkSlidingAtomic-8                11444576               100.8 ns/op             0 B/op          0 allocs/op
BenchmarkSlidingChannels
BenchmarkSlidingChannels-8               1966850               615.8 ns/op             0 B/op          0 allocs/op
BenchmarkSlidingConcurrentAtomic
BenchmarkSlidingConcurrentAtomic-8      27071841                41.05 ns/op            0 B/op          0 allocs/op
BenchmarkSlidingConcurrentChannels
BenchmarkSlidingConcurrentChannels-8     2029203               579.2 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/tetsuo/sliding        6.189s
```

## License

MIT License. See [LICENSE](LICENSE) for details.
