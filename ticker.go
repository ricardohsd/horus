package horus

import "time"

// TimeTicker defines a ticker interface for average counters
type TimeTicker interface {
	Stop()
	Chan() <-chan time.Time
}

// Ticker wraps time.Ticker to be TimeTicker complaint
type Ticker struct {
	*time.Ticker
}

var _ TimeTicker = (*Ticker)(nil)

func NewTicker(d time.Duration) *Ticker {
	return &Ticker{time.NewTicker(d)}
}

func (t *Ticker) Chan() <-chan time.Time {
	return t.C
}

// TestTicker provides control over a ticker. To be used only on tests.
type TestTicker struct {
	c chan time.Time
}

func NewTestTicker() *TestTicker {
	return &TestTicker{
		c: make(chan time.Time),
	}
}

func (t *TestTicker) Chan() <-chan time.Time {
	return t.c
}

func (t *TestTicker) Stop() {
}

func (t *TestTicker) Tick() {
	t.c <- time.Now()
	time.Sleep(100 * time.Millisecond)
}
