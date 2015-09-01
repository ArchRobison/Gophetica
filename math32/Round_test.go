package math32

import (
"testing"
"math"
)

func TestRound(t *testing.T) {
    ival := []float32{0} 
    x := float32(1)
    for !math.IsInf(float64(x),1) {
	    ival = append(ival,x)
	    x *= 2
    }
    fval := []float32{} 
    x = 0.5
    for x!=0 {
	    fval = append(fval,x,math.Nextafter32(x,1),math.Nextafter32(x,-1))
	    x /= 2
    }
    for _,i:=range ival {
	    for _,f:=range fval {
		    x := i+f
		    y := Round(x)
			d := Abs(y-x)
		    assert( d<=0.5, "Round/Abs", t)
		    assert( d<0.5 || int(y)&1==0, "Round/Abs", t)
        }
    }
}
