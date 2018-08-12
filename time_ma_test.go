package average

import (
	"math"
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
}

func TestTimeMA_Add(t *testing.T) {
	sma := &timeMA{
		window:      10 * time.Second,
		granularity: 2 * time.Second,
		size:        5,
		values:      make([]float64, 5),
	}

	tickers := make(chan time.Time)

	go sma.cleanBuckets(tickers)
	tickers <- time.Now()

	sma.Add(200.0)
	sma.Add(200.0)

	// Advance 2 positions
	tickers <- time.Now()
	tickers <- time.Now()
	time.Sleep(100 * time.Millisecond)

	sma.Add(10.0)
	sma.Add(20.0)

	tickers <- time.Now()
	time.Sleep(100 * time.Millisecond)

	sma.Add(25.0)

	tickers <- time.Now()
	time.Sleep(100 * time.Millisecond)

	sma.Add(5.0)
	sma.Add(15.0)

	tickers <- time.Now()
	time.Sleep(100 * time.Millisecond)

	sma.Add(25.0)

	avg := math.Round(sma.Average()*100) / 100
	total := 20.0
	if avg != total {
		t.Errorf("Average doesn't match. Expected %v, got %v", total, avg)
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
