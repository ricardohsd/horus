# Moving Averages

This package implements moving averages.

# Simple Moving Average

Accepts an int higher than zero as window size.

```go
sma, err := average.NewSMA(5)
if err != nil {
	panic(err)
}

sma.Add(11.0)
sma.Add(22.0)
sma.Add(33.0)
sma.Add(44.0)
sma.Add(55.0)
sma.Add(66.0)

fmt.Println(sma.Average()) #=> 44.0
```

# Cumulative Moving Average

TODO