package average

import (
	"math"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
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
}

func TestTimeMA_Add(t *testing.T) {
	clock := clockwork.NewFakeClock()
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make(map[int64]float64),
		clock:       clock,
	}

	sma.Add(10.0)
	sma.Add(20.0)

	clock.Advance(2 * time.Second)

	sma.Add(30.0)

	clock.Advance(4 * time.Second)

	sma.Add(20.0)
	sma.Add(10.0)

	clock.Advance(5 * time.Second)

	sma.Add(40.0)

	avg := math.Round(sma.Average()*100) / 100
	total := 20.0
	if avg != total {
		t.Errorf("Average doesn't match. Expected %v, got %v", total, avg)
	}
}

func TestTimeMA_Average(t *testing.T) {
	clock := clockwork.NewFakeClock()
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 1 * time.Second,
		values:      make(map[int64]float64),
		size:        10,
		clock:       clock,
	}

	for i := 0; i < 10; i++ {
		sma.Add(10.0)
		clock.Advance(1 * time.Second)
	}

	avg := sma.Average()
	total := 10.0
	if avg != total {
		t.Errorf("Average doesn't match. Expected %v, got %v", total, avg)
	}
}

func TestTimeMA_cleaning(t *testing.T) {
	clock := clockwork.NewRealClock()
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		values:      make(map[int64]float64),
		clock:       clock,
	}
	go sma.cleanBuckets()

	sma.Add(10.0)
	sma.Add(20.0)

	// wait for the window to be over
	clock.Sleep(12 * time.Second)

	avg := sma.Average()
	if avg != 0 {
		t.Errorf("Average should be 0, got %v", avg)
	}
}
