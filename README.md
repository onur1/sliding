# sliding

A Counter is a sliding window based counter for counting how frequently events are happening within some time window.

Note that, this implementation uses [ring](https://github.com/onur1/ring) under the hood, the ring capacity will always be a power of 2. For better precision, the given duration value needs to be rounded to the nearest millisecond which is divisible by `ring capacity - 1`.

[See the full API documentation at pkg.go.dev](https://pkg.go.dev/github.com/onur1/sliding)

## Example

```golang
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
```

Output:

```
0
2
3
1
```
