package average

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func TestTimeMA_Errors(t *testing.T) {
	_, err := NewTimeMA(0*time.Second, 1*time.Second)
	if err == nil {
		t.Errorf("Error can't be nil. Time SMA window must be higher than 0.")
	}

	_, err = NewTimeMA(1*time.Second, 0*time.Second)
	if err == nil {
		t.Errorf("Error can't be nil. Time SMA granularity must be higher than 0.")
	}

	_, err = NewTimeMA(1*time.Second, 2*time.Second)
	if err == nil {
		t.Errorf("Error can't be nil. Time SMA window must be a multiplier of granularity.")
	}
}

func TestTimeMA_Add(t *testing.T) {
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
	}

	ticker := NewTestTicker()

	go sma.cleanBuckets(ticker)
	ticker.Tick()

	sma.Add(200.0)
	sma.Add(200.0)

	// Tick to advance 2 positions
	ticker.Tick()
	ticker.Tick()

	sma.Add(10.0)
	sma.Add(20.0)

	ticker.Tick()

	sma.Add(25.0)

	ticker.Tick()

	sma.Add(5.0)
	sma.Add(15.0)

	ticker.Tick()

	sma.Add(25.0)

	avg := math.Round(sma.Average()*100) / 100
	total := 20.0
	if avg != total {
		t.Errorf("Average doesn't match. Expected %v, got %v", total, avg)
	}
}

func TestTimeMA_quit(t *testing.T) {
	sma := &timeMA{
		window:      5 * time.Second,
		granularity: 1 * time.Second,
		size:        5,
		values:      make([]float64, 5),
		quitC:       make(chan struct{}),
	}

	ticker := NewTestTicker()

	go sma.cleanBuckets(ticker)

	ticker.Tick()

	sma.Add(200.0)
	sma.Add(200.0)

	expected := []float64{
		0, 400, 0, 0, 0,
	}
	if !reflect.DeepEqual(expected, sma.values) {
		t.Errorf("Values don't match. Expected %v, got %v", expected, sma.values)
	}

	sma.Stop()

	sma.Add(500.0)

	if !reflect.DeepEqual(expected, sma.values) {
		t.Errorf("Values shouldn't be changed after stop. Expected %v, got %v", expected, sma.values)
	}
}

func TestTimeMA_cleaning(t *testing.T) {
	sma, err := NewTimeMA(5*time.Second, 1*time.Second)
	if err != nil {
		t.Errorf("Error should be nil. Got %v", err)
	}

	sma.Add(10.0)
	sma.Add(20.0)

	// wait for the window to be over
	time.Sleep(7 * time.Second)

	avg := sma.Average()
	if avg != 0 {
		t.Errorf("Average should be 0, got %v", avg)
	}
}
