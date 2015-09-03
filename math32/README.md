# math32

One of the annoyances of Go is that its standard math library does only float64 arithmetic.
Since Go has no implicit numerical conversions, using that library for 32-bit work is tedious.  

The library here also adds a function `Round` that rounds to the nearest integer,
with ties resolved as "round to even".

The library here corresponds to only a subset of the standard library.  It may grow over time.
