package sliding

import (
	"math"
	"sync/atomic"
	"time"
)

type slot struct {
	timestamp int64  // The timestamp of the slot in slot units
	count     uint64 // The count of events in this slot
}

// Counter is a lock-free sliding window counter implementation using atomic CAS operations.
type Counter struct {
	slots      []slot        // Ring buffer of slots
	totalSlots int64         // Total number of slots (power of 2 for efficient indexing)
	slotSize   time.Duration // Duration of each slot
	windowSize time.Duration // Total duration of the sliding window
	mask       int64         // Mask for efficient modulo operation
}

// NewCounter initializes a new sliding window counter.
//   - `windowSize` defines the total sliding window duration.
//   - The number of slots is calculated automatically for optimal performance.
func NewCounter(windowSize time.Duration) *Counter {
	// Convert window size to milliseconds for calculation
	ms := windowSize.Milliseconds()

	// Calculate the initial number of slots based on the logarithm of the window size
	slots := int(math.Ceil(math.Log(float64(ms))/math.Log(2))) + 1

	// Adjust slots to the next power of 2
	n := 1
	for n < slots {
		n <<= 1
	}
	slots = n

	// Recalculate window size to be an exact multiple of the number of slots
	slotSize := windowSize / time.Duration(slots)
	windowSize = slotSize * time.Duration(slots)

	// Prepare the ring buffer mask for efficient indexing
	mask := int64(slots - 1)

	slotsArray := make([]slot, slots)
	return &Counter{
		slots:      slotsArray,
		totalSlots: int64(slots),
		slotSize:   slotSize,
		windowSize: windowSize,
		mask:       mask,
	}
}

// Increment increments the count in the current time slot.
func (c *Counter) Increment() {
	now := time.Now()
	slotTime := now.UnixNano() / c.slotSize.Nanoseconds()
	idx := slotTime & c.mask

	for {
		slotPtr := &c.slots[idx]
		slotTimestamp := atomic.LoadInt64(&slotPtr.timestamp)

		if slotTimestamp == slotTime {
			// Current slot; increment the count atomically
			atomic.AddUint64(&slotPtr.count, 1)
			return
		} else if slotTimestamp < slotTime {
			// Outdated slot; attempt to reset it.
			if atomic.CompareAndSwapInt64(&slotPtr.timestamp, slotTimestamp, slotTime) {
				// Successfully updated timestamp; reset count to 1
				atomic.StoreUint64(&slotPtr.count, 1)
				return
			}
			// CAS failed; another goroutine updated the slot, retry
		} else {
			// Slot timestamp is ahead (possible time skew); increment count to be safe
			atomic.AddUint64(&slotPtr.count, 1)
			return
		}
	}
}

// Peek returns the total count across the sliding window.
func (c *Counter) Peek() uint64 {
	now := time.Now()
	slotTime := now.UnixNano() / c.slotSize.Nanoseconds()
	windowStart := slotTime - c.totalSlots + 1

	var total uint64 = 0
	for i := int64(0); i < c.totalSlots; i++ {
		idx := (slotTime - i) & c.mask
		slotPtr := &c.slots[idx]
		slotTimestamp := atomic.LoadInt64(&slotPtr.timestamp)

		if slotTimestamp >= windowStart {
			count := atomic.LoadUint64(&slotPtr.count)
			total += count
		}
	}
	return total
}

// FrameDuration returns the duration of each slot (frame) in the sliding window.
func (c *Counter) FrameDuration() time.Duration {
	return c.slotSize
}
