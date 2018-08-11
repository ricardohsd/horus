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

# Moving Average, time window

Calculates the average of values in a given moving window.

```go
ma, err := average.NewTimeMA(10 * time.Duration, 2 * time.Duration)
if err != nil {
	panic(err)
}

ma.Add(11.0)
ma.Add(22.0)

time.Sleep(11 * time.Second)

ma.Add(33.0)
ma.Add(44.0)
ma.Add(55.0)
ma.Add(66.0)

fmt.Println(ma.Average()) #=> 19.80
```