package horus

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrWrongWindowSize = errors.New("specified window must be equal or less than total window")

type RollingWindow struct {
	sync.RWMutex
	window      time.Duration
	granularity time.Duration
	size        int
	position    int
	now         time.Time
	values      []float64
	counters    []int64
	quitC       chan struct{}
	close       bool
}

// NewRWindow provides a slide window operation for a moving average in a rolling window.
// The given window must be in a valid time.Duration higher than 0.
func NewRWindow(window time.Duration, granularity time.Duration) (*RollingWindow, error) {
	if window == 0 {
		return nil, fmt.Errorf("window must be higher than 0")
	}

	if granularity == 0 {
		return nil, fmt.Errorf("granularity must be higher than zero")
	}

	if granularity > window || window%granularity != 0 {
		return nil, fmt.Errorf("window must be a multiplier of granularity")
	}

	size := int(window / granularity)

	t := &RollingWindow{
		window:      window,
		granularity: granularity,
		size:        size,
		position:    0,
		values:      make([]float64, size),
		counters:    make([]int64, size),
		quitC:       make(chan struct{}),
		close:       false,
	}

	ticker := NewTicker(t.granularity)

	go t.cleanBuckets(ticker)

	return t, nil
}

func (r *RollingWindow) cleanBuckets(ticker TimeTicker) {
	for {
		select {
		case tick := <-ticker.Chan():
			r.Lock()

			r.position++

			if r.position >= r.size {
				r.position = 0
			}

			r.values[r.position] = 0
			r.counters[r.position] = 0

			r.now = tick

			r.Unlock()
		case <-r.quitC:
			ticker.Stop()
			return
		}
	}
}

func (r *RollingWindow) Stop() {
	r.Lock()
	defer r.Unlock()

	r.close = true
	r.quitC <- struct{}{}
}

// Add adds a value to its given position in the rolling window.
// The given value is incremented with the position's value.
func (r *RollingWindow) Add(value float64) {
	r.Lock()
	defer r.Unlock()

	if r.close {
		return
	}

	r.counters[r.position]++
	r.values[r.position] += value
}

// AddWithTime takes a value and timestamp and adds it into the correct position.
// The operation is ignored if the timestamp is older than the window size.
func (r *RollingWindow) AddWithTime(value float64, t time.Time) {
	r.Lock()
	defer r.Unlock()

	if r.close {
		return
	}

	// Compare transaction timestamp with current time - window in seconds
	// if it is older, exit function.
	if t.Unix() < r.now.Add(-r.window).Unix() {
		return
	}

	diff := int(r.now.Unix() - t.Unix())
	tpos := r.position - diff
	if tpos < 0 {
		tpos = r.size + tpos
	}

	r.counters[tpos]++
	r.values[tpos] += value
}

// Count returns the total number of transactions in the rolling window
func (r *RollingWindow) Count() int64 {
	r.RLock()
	defer r.RUnlock()

	var total int64
	total = 0

	for _, v := range r.counters {
		total += v
	}

	return total
}

// Average calculates the average inside rolling window.
func (r *RollingWindow) Average() float64 {
	r.RLock()
	defer r.RUnlock()

	total := 0.0

	for _, v := range r.values {
		total = total + v
	}

	if total == 0 {
		return total
	}

	return total / float64(r.size)
}

// AverageSince calculates the average in the given window
func (r *RollingWindow) AverageSince(w time.Duration) (float64, error) {
	r.RLock()
	defer r.RUnlock()

	if w > r.window {
		return 0, ErrWrongWindowSize
	}

	sum := 0.0
	count := 0.0
	windowSize := int(w / r.granularity)

	for i := 0; i < windowSize; i++ {
		pos := r.position - i
		if pos < 0 {
			pos += len(r.values)
		}

		sum += r.values[pos]
		count++
	}

	return sum / count, nil
}

// Max returns the max value in the given rolling window.
func (r *RollingWindow) Max() float64 {
	r.RLock()
	defer r.RUnlock()

	max := r.values[0]

	for _, v := range r.values {
		if v > max {
			max = v
		}
	}

	return max
}

// Min returns the min value in the given rolling window.
func (r *RollingWindow) Min() float64 {
	r.RLock()
	defer r.RUnlock()

	min := r.values[0]

	for _, v := range r.values {
		if v <= min {
			min = v
		}
	}

	return min
}
