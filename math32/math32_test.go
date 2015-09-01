package math32

import (
	"testing"
	"math"
)

func assert(expr bool, msg string, t *testing.T) {
	if !expr {
		t.Fail()
		panic(msg)
	}
}

func TestMath32(t *testing.T) {
	val := []float32{-2.25, -1.0, -2.25, -0.0, 0.0, 1.0, 2.25}
	for _, x := range val {
		w := Abs(x)
		assert(w == Max(x, -x), "Abs/Max", t)
		assert(-w == Min(x, -x), "Abs/Min", t)
		if x >= 0 {
			z := Sqrt(x)
			assert(z*z == x, "Sqrt", t)
		}
		for _, y := range val {
			a := Min(x, y)
			b := Max(x, y)
			assert(a <= b, "Min/Max", t)
			assert(x == y || a < b, "Min/Max", t)
			assert(a == x || a == y, "Min", t)
			assert(b == x || b == y, "Max", t)
			θ := Atan2(y, x)
			r := Hypot(x, y)
			v, u := Sincos(θ)
			assert(Sin(θ) == v, "SinCos/Sin", t)
			assert(Cos(θ) == u, "SinCos/Cos", t)
			assert(Abs(u*r-x) < 0.000001, "Atan2/Hypot/Cos", t)
			assert(Abs(v*r-y) < 0.000001, "Atan2/Hypot/Sin", t)
		}
	}
}

func TestTrunc(t *testing.T) {
    ival := []float32{} 
    x := float32(1)
    for !math.IsInf(float64(x),1) {
	    ival = append(ival,x)
	    x *= 2
    }
    fval := []float32{} 
    x = math.Nextafter32(1,-1)
    for x!=0 {
	    fval = append(fval,x)
	    x = 1 - (1-x)*2
    }
    for _,i:=range ival {
	    for _,f:=range fval {
		    x := i+f
	        y := Trunc(x)
			assert( y==i || x-i==1, "Trunc", t)
		    z := Trunc(-x)
			assert( z==-y, "Trunc", t)
		}
    }
}

