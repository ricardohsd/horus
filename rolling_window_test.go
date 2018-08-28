package horus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRollingWindow_Errors(t *testing.T) {
	_, err := NewRWindow(0*time.Second, 1*time.Second)
	assert.NotNil(t, err, "Error can't be nil. Window must be higher than 0.")

	_, err = NewRWindow(1*time.Second, 0*time.Second)
	assert.NotNil(t, err, "Error can't be nil. Granularity must be higher than 0.")

	_, err = NewRWindow(1*time.Second, 2*time.Second)
	assert.NotNil(t, err, "Error can't be nil. Window must be a multiplier of granularity.")
}

func TestRollingWindow(t *testing.T) {
	rw := &RollingWindow{
		window:      10 * time.Second,
		granularity: time.Second,
		size:        10,
		values:      make([]float64, 10),
		counters:    make([]int64, 10),
	}

	ticker := NewTestTicker()

	go rw.cleanBuckets(ticker)
	ticker.Tick()

	assert.Equal(t, 0.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())

	rw.Add(200.0)
	rw.Add(200.0)

	assert.Equal(t, 400.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())
	assert.Equal(t, 40.0, rw.Average())

	// Tick to advance 6 positions
	for i := 0; i < 6; i++ {
		ticker.Tick()
	}

	rw.Add(10.0)
	rw.Add(20.0)

	ticker.Tick()

	rw.Add(25.0)

	ticker.Tick()

	rw.Add(5.0)
	rw.Add(15.0)

	ticker.Tick()

	rw.Add(-25.0)

	ticker.Tick()

	rw.Add(80.0)

	assert.Equal(t, 80.0, rw.Max())
	assert.Equal(t, -25.0, rw.Min())
	assert.Equal(t, 13.0, rw.Average())
	assert.Equal(t, int64(7), rw.Count())
}

func TestRollingWindow_WithTime(t *testing.T) {
	rw := &RollingWindow{
		window:      10 * time.Second,
		granularity: time.Second,
		size:        10,
		values:      make([]float64, 10),
		counters:    make([]int64, 10),
	}

	ticker := NewTestTicker()

	go rw.cleanBuckets(ticker)
	now := ticker.Tick()

	assert.Equal(t, 0.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())

	rw.AddWithTime(200.0, now)
	rw.AddWithTime(200.0, now)

	assert.Equal(t, 400.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())
	assert.Equal(t, 40.0, rw.Average())

	rw.AddWithTime(150.0, now.Add(-5*time.Second))

	assert.Equal(t, 400.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())
	assert.Equal(t, 55.0, rw.Average())

	rw.AddWithTime(10000.0, now.Add(-15*time.Second))

	assert.Equal(t, 400.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())
	assert.Equal(t, 55.0, rw.Average())

	rw.AddWithTime(-10.0, now.Add(-10*time.Second))

	assert.Equal(t, 390.0, rw.Max())
	assert.Equal(t, 0.0, rw.Min())
	assert.Equal(t, 54.0, rw.Average())
}
func TestRollingWindow_AverageSince(t *testing.T) {
	rw := &RollingWindow{
		window:      6 * time.Second,
		granularity: time.Second,
		size:        6,
		values: []float64{
			0, 10, -20, 120, 10, 30,
		},
		counters: make([]int64, 5),
	}

	rw.position = 5

	_, err := rw.AverageSince(1 * time.Hour)
	assert.NotNil(t, err)

	avg, err := rw.AverageSince(6 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 25.0, avg)

	avg, err = rw.AverageSince(3 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 53.333333333333336, avg)

	avg, err = rw.AverageSince(5 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 30.0, avg)

	rw.position = 3

	avg, err = rw.AverageSince(6 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 25.0, avg)

	avg, err = rw.AverageSince(2 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 50.0, avg)

	rw.position = 0

	avg, err = rw.AverageSince(6 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 25.0, avg)

	avg, err = rw.AverageSince(2 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, 15.0, avg)
}

func TestRollingWindow_MaxMin(t *testing.T) {
	rw := &RollingWindow{
		values: []float64{
			0, 10, -20, 150, 0, 30,
		},
	}

	assert.Equal(t, 150.0, rw.Max())
	assert.Equal(t, -20.0, rw.Min())
}

func TestRollingWindow_quit(t *testing.T) {
	rw := &RollingWindow{
		window:      5 * time.Second,
		granularity: 1 * time.Second,
		size:        5,
		values:      make([]float64, 5),
		counters:    make([]int64, 5),
		quitC:       make(chan struct{}),
	}

	ticker := NewTestTicker()

	go rw.cleanBuckets(ticker)

	ticker.Tick()

	rw.Add(200.0)
	rw.Add(200.0)

	expected := []float64{
		0, 400, 0, 0, 0,
	}
	assert.Equal(t, expected, rw.values)

	rw.Stop()

	rw.Add(500.0)

	assert.Equal(t, expected, rw.values)
}

func TestRollingWindow_cleaning(t *testing.T) {
	rw, err := NewRWindow(5*time.Second, 1*time.Second)
	assert.Nil(t, err)

	rw.Add(10.0)
	rw.Add(20.0)

	// wait for the window to be over
	time.Sleep(7 * time.Second)

	assert.Equal(t, 0.0, rw.Average())
}
