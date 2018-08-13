package horus

import (
	"math"
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

func TestRollingWindow_Add(t *testing.T) {
	rw := &rollingWindow{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
		counters:    make([]int, 5),
	}

	ticker := NewTestTicker()

	go rw.cleanBuckets(ticker)
	ticker.Tick()

	rw.Add(200.0)
	rw.Add(200.0)

	// Tick to advance 2 positions
	ticker.Tick()
	ticker.Tick()

	rw.Add(10.0)
	rw.Add(20.0)

	ticker.Tick()

	rw.Add(25.0)

	ticker.Tick()

	rw.Add(5.0)
	rw.Add(15.0)

	ticker.Tick()

	rw.Add(25.0)

	avg := math.Round(rw.Average()*100) / 100
	assert.Equal(t, 20.0, avg)

	assert.Equal(t, 6, rw.Count())
}

func TestRollingWindow_MaxMin(t *testing.T) {
	rw := &rollingWindow{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
		counters:    make([]int, 5),
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

	// Tick to advance 2 positions
	ticker.Tick()
	ticker.Tick()

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

	assert.Equal(t, 7, rw.Count())
}

func TestRollingWindow_quit(t *testing.T) {
	rw := &rollingWindow{
		window:      5 * time.Second,
		granularity: 1 * time.Second,
		size:        5,
		values:      make([]float64, 5),
		counters:    make([]int, 5),
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
