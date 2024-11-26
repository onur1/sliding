
# sliding

**sliding** is a lock-free sliding window counter implementation using atomic CAS operations.

## Installation

Install the library by running:

```sh
go get github.com/onur1/sliding
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

```
goos: darwin
goarch: amd64
pkg: github.com/onur1/sliding
cpu: Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz
BenchmarkSlidingAtomic
    sliding_benchmark_test.go:20: expected=1 result=1 diff=0
    sliding_benchmark_test.go:20: expected=100 result=100 diff=0
    sliding_benchmark_test.go:20: expected=10000 result=10000 diff=0
    sliding_benchmark_test.go:20: expected=1000000 result=1000000 diff=0
    sliding_benchmark_test.go:20: expected=11295186 result=9360845 diff=1934341
BenchmarkSlidingAtomic-8                11295186               101.6 ns/op             0 B/op          0 allocs/op
BenchmarkSlidingChannels
    sliding_benchmark_test.go:31: expected=1 result=1 diff=0
    sliding_benchmark_test.go:31: expected=100 result=100 diff=0
    sliding_benchmark_test.go:31: expected=10000 result=10000 diff=0
    sliding_benchmark_test.go:31: expected=1000000 result=1000000 diff=0
    sliding_benchmark_test.go:31: expected=1894435 result=1574559 diff=319876
BenchmarkSlidingChannels-8               1894435               630.4 ns/op             0 B/op          0 allocs/op
BenchmarkSlidingConcurrentAtomic
    sliding_benchmark_test.go:54: expected=100 result=100 diff=0
    sliding_benchmark_test.go:54: expected=10000 result=10000 diff=0
    sliding_benchmark_test.go:54: expected=1000000 result=1000000 diff=0
    sliding_benchmark_test.go:54: expected=28119298 result=21187932 diff=6931366
BenchmarkSlidingConcurrentAtomic-8      28119298                44.15 ns/op            0 B/op          0 allocs/op
BenchmarkSlidingConcurrentChannels
    sliding_benchmark_test.go:77: expected=100 result=100 diff=0
    sliding_benchmark_test.go:77: expected=10000 result=10000 diff=0
    sliding_benchmark_test.go:77: expected=1000000 result=1000000 diff=0
    sliding_benchmark_test.go:77: expected=1958114 result=1628152 diff=329962
BenchmarkSlidingConcurrentChannels-8     1958114               602.0 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/onur1/sliding        6.737s
```

## License

MIT License. See [LICENSE](LICENSE) for details.
