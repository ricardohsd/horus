# Horus

This package implements statistics in a fixed or time rolling window.

# Simple Moving window

Accepts an int higher than zero as window size.

```go
sw, err := horus.NewSWindow(5)
if err != nil {
	panic(err)
}

sw.Add(11.0)
sw.Add(22.0)
sw.Add(33.0)
sw.Add(44.0)
sw.Add(55.0)
sw.Add(66.0)

fmt.Println(sw.Average()) #=> 44.0
fmt.Println(sw.Min()) #=> 22.0
fmt.Println(sw.Max()) #=> 66.0
```

# Time Rolling window

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