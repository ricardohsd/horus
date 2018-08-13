package horus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMA_ZeroWindow(t *testing.T) {
	_, err := NewSMA(0)
	assert.NotNil(t, err, "Error shouldn't be nil. Average window must be higher than zero.")
}

func TestSMA(t *testing.T) {
	sma, err := NewSMA(5)
	assert.Nil(t, err)

	sma.Add(10.0)
	sma.Add(5.0)
	sma.Add(15.0)
	sma.Add(16.0)
	sma.Add(4.0)

	assert.Equal(t, 10.0, sma.Average())
	assert.Equal(t, 4.0, sma.Min())
	assert.Equal(t, 16.0, sma.Max())

	sma.Add(20.0)

	assert.Equal(t, 12.0, sma.Average())
	assert.Equal(t, 4.0, sma.Min())
	assert.Equal(t, 20.0, sma.Max())

	sma.Add(-13.0)

	assert.Equal(t, 8.4, sma.Average())
	assert.Equal(t, -13.0, sma.Min())
	assert.Equal(t, 20.0, sma.Max())
}
