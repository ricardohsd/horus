# Horus

This package implements statistics in a fixed or rolling window.

# Simple Moving Average

Accepts an int higher than zero as window size.

```go
sma, err := horus.NewSMA(5)
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

# Rolling window

```go
rw, err := horus.NewRWindow(10*time.Second, 2*time.Second)
if err != nil {
	panic(err)
}

rw.Add(11.0)
rw.Add(22.0)

time.Sleep(11 * time.Second)

rw.Add(33.0)
rw.Add(44.0)
time.Sleep(1 * time.Second)
rw.Add(55.0)
rw.Add(66.0)

rw.Stop()

fmt.Println(rw.Average()) #=> 39.6
fmt.Println(rw.Max()) #=> 121
fmt.Println(rw.Min()) #=> 0
```