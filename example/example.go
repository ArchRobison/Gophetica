package main

import (
	"math"
    nimble "github.com/ArchRobison/NimbleDraw"
)

type Banana struct {
	value int32
}

func (b *Banana) Init(width, height int32) {
	b.value = 0
}

var w = make([]float32, 44100)

var c = nimble.RGB(0,0,0)

func init() {
		a := float32(1)
		for i := range w {
			w[i] = float32(math.Sin(float64(i) * (2*math.Pi*440/nimble.SampleRate) ))*a
		    a *= 0.9998
		}
}

func (b *Banana) Render(pm nimble.PixMap) {
	if b.value == 0 {
		nimble.PlaySound(w[:22500], 1.0, 0.5/2+.01)
		c = nimble.RGB(1,0,0)
    }
    if b.value == 60 {
		nimble.PlaySound(w[:34500], 1.0, .75/2+.01)
		c = nimble.RGB(0,1,0)
    }
	if b.value == 120 {
		nimble.PlaySound(w, 1.0, 1.0/2+.01)
		c = nimble.RGB(0,0,1)
	}
	pm.DrawRect(nimble.Rect{Left: 10, Top: 20, Right: 100 + b.value, Bottom: 50 + b.value}, c)
	b.value++
}

func main() {
	nimble.AddRenderClient(&Banana{})
	nimble.Run()
}
