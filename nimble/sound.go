package nimble

import "fmt"

const SampleRate = 44100.0

func PlaySound(waveform []float32, relativeAmplitude, relativePitch float32) {
	select {
	case newPlayers <- player{
		waveform:  waveform,
		delta:     fixedPointTime(relativePitch*fixedPointUnit + 0.5),
		amplitude: relativeAmplitude,
		index:     0,
	}:
	default:
		// Channel full, so drop sound.
		if devConfig {
			fmt.Printf("PlaySound: channel full\n")
		}
	}
}

const (
	fixedPointShift = 10
	fixedPointUnit  = 1 << fixedPointShift
	fixedPointMask  = fixedPointUnit - 1
	fixedPointScale = 1. / fixedPointUnit
)

type fixedPointTime uint32

type player struct {
	waveform  []float32
	amplitude float32
	delta     fixedPointTime
	index     fixedPointTime // Index into waveform
}

var newPlayers chan player = make(chan player, 100)

var activePlayers []player

// Assumes buf has already been zeroed.
func getSoundSamples(buf []float32) {
	// Get any new players
drain:
	for {
		select {
		case p := <-newPlayers:
			activePlayers = append(activePlayers, p)
		default:
			break drain
		}
	}
	j := len(activePlayers)
playerLoop:
	for i := 0; i < j; {
		p := &activePlayers[i]
		for k := range buf {
			m := int(p.index >> fixedPointShift)
			if m+1 >= len(p.waveform) {
				// Reached end of waveform.  Delete the player.
				j--
				*p = activePlayers[j]
				continue playerLoop
			}
			// Interpolate
			f := float32(p.index&fixedPointMask) * fixedPointScale
			w := p.waveform[m]*(1-f) + p.waveform[m+1]*f
			buf[k] += w * p.amplitude
			p.index += p.delta
		}
		i++
	}
	activePlayers = activePlayers[:j]
}
