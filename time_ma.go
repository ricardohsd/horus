package average

import (
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

type timeMA struct {
	sync.RWMutex
	window      time.Duration
	granularity time.Duration
	size        int
	position    int
	values      []float64
	clock       clockwork.Clock
}

// NewTimeMA provides a slide window operation for a moving average in a time window.
// The given window must be in a valid time.Duration higher than 0.
func NewTimeMA(window time.Duration, granularity time.Duration) (*timeMA, error) {
	if window == 0 {
		return nil, fmt.Errorf("window must be higher than 0")
	}

	if granularity == 0 {
		return nil, fmt.Errorf("granularity must be higher than zero")
	}

	t := &timeMA{
		window:      window,
		granularity: granularity,
		size:        int(window / granularity),
		position:    0,
		values:      make([]float64, int(window/granularity)),
		clock:       clockwork.NewRealClock(),
	}

	go t.cleanBuckets()

	return t, nil
}

func (t *timeMA) cleanBuckets() {
	ticker := t.clock.NewTicker(t.granularity)

	for {
		select {
		case <-ticker.Chan():
			t.Lock()

			t.position++

			if t.position >= t.size {
				t.position = 0
			}

			t.values[t.position] = 0

			t.Unlock()
		}
	}
}

// Add adds a value to its given bucket in the time window.
// Each value is added as a separate value in the internal storage.
func (t *timeMA) Add(value float64) {
	t.Lock()
	defer t.Unlock()

	t.values[t.position] += value
}

// Average calculates the average of buckets inside time window.
func (t *timeMA) Average() float64 {
	t.RLock()
	defer t.RUnlock()

	total := float64(0)

	for _, v := range t.values {
		total = total + v
	}

	if total == 0 {
		return total
	}

	return total / float64(t.size)
}
