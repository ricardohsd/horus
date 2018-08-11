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
	values      map[int64]float64
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
		values:      make(map[int64]float64),
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

			lb := t.clock.Now().Add(-t.window).Unix()

			for k := range t.values {
				if k < lb {
					delete(t.values, k)
				}
			}

			t.Unlock()
		}
	}
}

// Add adds a value to its given bucket in the time window.
// Each value is added as a separate value in the internal storage.
func (t *timeMA) Add(value float64) {
	t.Lock()
	defer t.Unlock()

	now := t.clock.Now()
	bucket := now.Round(t.granularity).Unix()

	v, ok := t.values[bucket]
	if !ok {
		v = value
	} else {
		v = value + v
	}

	t.values[bucket] = v
}

// Average calculates the average of buckets inside time window.
func (t *timeMA) Average() float64 {
	t.RLock()
	defer t.RUnlock()

	lb := t.clock.Now().Add(-t.window).Unix()

	total := float64(0)

	for k, v := range t.values {
		if k >= lb {
			total = total + v
		}
	}

	if total == 0 {
		return total
	}

	return total / float64(t.size)
}
