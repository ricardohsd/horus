package average

import (
	"testing"
)

func TestSMA_ZeroWindow(t *testing.T) {
	_, err := NewSMA(0)
	if err == nil {
		t.Errorf("Error shouldn't be nil. Average window must be higher than zero.")
	}
}

func TestSMA(t *testing.T) {
	sma, err := NewSMA(5)
	if err != nil {
		t.Errorf("Error should be nil. Got %v", err)
	}

	sma.Add(10.0)
	sma.Add(10.0)
	sma.Add(10.0)
	sma.Add(10.0)
	sma.Add(10.0)

	avg := sma.Average()
	if avg != 10.0 {
		t.Errorf("Average doesn't match. Expected %v, got %v", 10.0, avg)
	}

	sma.Add(20.0)

	avg = sma.Average()
	if avg != 12.0 {
		t.Errorf("Average doesn't match. Expected %v, got %v", 12.0, avg)
	}

	sma.Add(-13.0)

	avg = sma.Average()
	if avg != 7.4 {
		t.Errorf("Average doesn't match. Expected %v, got %v", 7.4, avg)
	}
}
