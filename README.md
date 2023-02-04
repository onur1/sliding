# sliding

A Counter is a sliding window based counter for counting how frequently events are happening within some time window.

Note that, this implementation uses [ring](https://github.com/onur1/ring) under the hood, the ring capacity will always be a power of 2. For better precision, the given duration value needs to be rounded to the nearest millisecond which is divisible by `ring capacity - 1`.
