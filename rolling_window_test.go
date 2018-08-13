package horus

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestRollingWindow_Errors(t *testing.T) {
	_, err := NewRWindow(0*time.Second, 1*time.Second)
	if err == nil {
		t.Errorf("Error can't be nil. Wndow must be higher than 0.")
	}

	_, err = NewRWindow(1*time.Second, 0*time.Second)
	if err == nil {
		t.Errorf("Error can't be nil. Granularity must be higher than 0.")
	}

	_, err = NewRWindow(1*time.Second, 2*time.Second)
	if err == nil {
		t.Errorf("Error can't be nil. Window must be a multiplier of granularity.")
	}
}

func TestRollingWindow_Add(t *testing.T) {
	rw := &rollingWindow{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
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
	total := 20.0
	if avg != total {
		t.Errorf("Average doesn't match. Expected %v, got %v", total, avg)
	}
}

func TestRollingWindow_MaxMin(t *testing.T) {
	rw := &rollingWindow{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
	}

	ticker := NewTestTicker()

	go rw.cleanBuckets(ticker)
	ticker.Tick()

	max := rw.Max()
	expected := 0.0
	if max != expected {
		t.Errorf("Max doesn't match. Expected %v, got %v", expected, max)
	}

	min := rw.Min()
	expected = 0.0
	if min != expected {
		t.Errorf("Min doesn't match. Expected %v, got %v", expected, min)
	}

	rw.Add(200.0)
	rw.Add(200.0)

	max = rw.Max()
	expected = 400.0
	if max != expected {
		t.Errorf("Max doesn't match. Expected %v, got %v", expected, max)
	}

	min = rw.Min()
	expected = 0.0
	if min != expected {
		t.Errorf("Min doesn't match. Expected %v, got %v", expected, min)
	}

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

	max = rw.Max()
	expected = 80.0
	if max != expected {
		t.Errorf("Max doesn't match. Expected %v, got %v", expected, max)
	}

	min = rw.Min()
	expected = -25.0
	if min != expected {
		t.Errorf("Min doesn't match. Expected %v, got %v", expected, min)
	}
}

func TestRollingWindow_quit(t *testing.T) {
	rw := &rollingWindow{
		window:      5 * time.Second,
		granularity: 1 * time.Second,
		size:        5,
		values:      make([]float64, 5),
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
	if !reflect.DeepEqual(expected, rw.values) {
		t.Errorf("Values don't match. Expected %v, got %v", expected, rw.values)
	}

	rw.Stop()

	rw.Add(500.0)

	if !reflect.DeepEqual(expected, rw.values) {
		t.Errorf("Values shouldn't be changed after stop. Expected %v, got %v", expected, rw.values)
	}
}

func TestRollingWindow_cleaning(t *testing.T) {
	rw, err := NewRWindow(5*time.Second, 1*time.Second)
	if err != nil {
		t.Errorf("Error should be nil. Got %v", err)
	}

	rw.Add(10.0)
	rw.Add(20.0)

	// wait for the window to be over
	time.Sleep(7 * time.Second)

	avg := rw.Average()
	if avg != 0 {
		t.Errorf("Average should be 0, got %v", avg)
	}
}
