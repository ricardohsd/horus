package horus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMA_ZeroWindow(t *testing.T) {
	_, err := NewSWindow(0)
	assert.NotNil(t, err, "Error shouldn't be nil. Average window must be higher than zero.")
}

func TestSMA(t *testing.T) {
	sw, err := NewSWindow(5)
	assert.Nil(t, err)

	sw.Add(10.0)
	sw.Add(5.0)
	sw.Add(15.0)
	sw.Add(16.0)
	sw.Add(4.0)

	assert.Equal(t, 10.0, sw.Average())
	assert.Equal(t, 4.0, sw.Min())
	assert.Equal(t, 16.0, sw.Max())

	sw.Add(20.0)

	assert.Equal(t, 12.0, sw.Average())
	assert.Equal(t, 4.0, sw.Min())
	assert.Equal(t, 20.0, sw.Max())

	sw.Add(-13.0)

	assert.Equal(t, 8.4, sw.Average())
	assert.Equal(t, -13.0, sw.Min())
	assert.Equal(t, 20.0, sw.Max())
}
