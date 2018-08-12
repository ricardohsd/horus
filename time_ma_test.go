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
	clock := clockwork.NewRealClock()
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
		clock:       clock,
	}
	go sma.cleanBuckets()

	sma.Add(10.0)
	sma.Add(20.0)

	clock.Sleep(2 * time.Second)

	sma.Add(30.0)

	clock.Sleep(4 * time.Second)

	sma.Add(20.0)
	sma.Add(10.0)

	clock.Sleep(5 * time.Second)

	sma.Add(40.0)

	avg := math.Round(sma.Average()*100) / 100
	total := 20.0
	if avg != total {
		t.Errorf("Average doesn't match. Expected %v, got %v", total, avg)
	}
}

func TestTimeMA_cleaning(t *testing.T) {
	clock := clockwork.NewRealClock()
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		values:      make([]float64, 5),
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
