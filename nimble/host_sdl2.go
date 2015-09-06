// Plaform-dependent routines

package nimble

// typedef unsigned char Uint8;
// void getSoundSamplesAdaptor(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"reflect"
	"runtime"
	"unsafe"
)

type renderClient interface {
	Init(width, height int32) // Inform client of window size
	Render(pm PixMap)
}

var renderClientList []renderClient

func AddRenderClient(r renderClient) {
	renderClientList = append(renderClientList, r)
}

type keyClient interface {
	KeyDown(Key)
}

var keyClientList []keyClient

func AddKeyClient(k keyClient) {
	keyClientList = append(keyClientList, k)
}

var keyMap = map[sdl.Keycode]Key{
	sdl.K_RETURN:    KeyReturn,
	sdl.K_ESCAPE:    KeyEscape,
	sdl.K_LEFT:      KeyLeft,
	sdl.K_RIGHT:     KeyRight,
	sdl.K_UP:        KeyUp,
	sdl.K_DOWN:      KeyDown,
	sdl.K_DELETE:    KeyDelete,
	sdl.K_BACKSPACE: KeyBackspace,
	sdl.K_TAB:       KeyTab,
}

var mouseX, mouseY int32

// Get position of mouse
func MouseWhere() (x, y int32) {
	x = int32(mouseX)
	y = int32(mouseY)
	return
}

// Get time in seconds.  Time zero is platform specific.
func Now() float64 {
	return float64(sdl.GetTicks()) * 0.001
}

// Creates a slice of Pixel from a raw pointer
func sliceFromPixelPtr(data unsafe.Pointer, length int) (pixels []Pixel) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&pixels))
	sliceHeader.Cap = int(length)
	sliceHeader.Len = int(length)
	sliceHeader.Data = uintptr(data)
	return
}

func lockTexture(tex *sdl.Texture, width int, height int) (pixels []Pixel, pitch int) {
	var data unsafe.Pointer
	err := tex.Lock(nil, &data, &pitch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tex.Lock: %v", err)
		panic(err)
	}
	// Convert pitch units from byte to pixels
	pitch /= 4
	pixels = sliceFromPixelPtr(data, width*height)
	return
}

var winTitle string = "FIXME"
var winWidth, winHeight int = 800, 600

func sliceFromAudioStream(data unsafe.Pointer, length int) (samples []float32) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&samples))
	sliceHeader.Cap = int(length)
	sliceHeader.Len = int(length)
	sliceHeader.Data = uintptr(data)
	return
}

//export getSoundSamplesAdaptor
func getSoundSamplesAdaptor(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	buf := sliceFromAudioStream(unsafe.Pointer(stream), int(length)/4)
	for i := range buf {
		buf[i] = 0
	}
	getSoundSamples(buf)
}

func Run() int {
	// All SDL calls must come from same thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Install audio callback
	spec := &sdl.AudioSpec{
		Freq:     SampleRate,
		Format:   sdl.AUDIO_F32SYS,
		Channels: 1,
		Samples:  4096,
		Callback: sdl.AudioCallback(C.getSoundSamplesAdaptor),
	}
	audioDevice, err := sdl.OpenAudioDevice("", false, spec, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open audio device: %v\n", err)
		panic(err)
	}
	if audioDevice < 2 {
		fmt.Fprintf(os.Stderr, "Audio device=%v < 2 contrary to SDL-2 documentation\n", audioDevice, err)
	}
	sdl.PauseAudioDevice(audioDevice, false)

	// Create window
	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %v\n", err)
		panic(err)
	}
	defer window.Destroy()

	// Create renderer
	width, height := window.GetSize()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %v\n", err)
		panic(err)
	}
	defer renderer.Destroy()

	// Create texture
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		fmt.Fprintf(os.Stderr, "renderer.CreateTexture: %v\n", err)
		panic(err)
	}
	defer tex.Destroy()

	for _, r := range renderClientList {
		r.Init(int32(width), int32(height))
	}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return 0
			case *sdl.MouseMotionEvent:
				mouseX = int32(e.X)
				mouseY = int32(e.Y)
			case *sdl.KeyDownEvent:
				var k Key
				if 0x20 <= e.Keysym.Sym && e.Keysym.Sym < 0x7F {
					// Printable ASCII
					k = Key(e.Keysym.Sym)
				} else {
					// Try special character table
					k = keyMap[e.Keysym.Sym]
				}
				if k != 0 {
					for _, c := range keyClientList {
						c.KeyDown(k)
					}
				}
			}
		}

		pixels, pitch := lockTexture(tex, width, height)
		pm := MakePixMap(int32(width), int32(height), pixels, int32(pitch))
		for _, r := range renderClientList {
			r.Render(pm)
		}
		tex.Unlock()

		err := renderer.Clear()
		if err != nil {
			fmt.Fprintf(os.Stderr, "renderer.Clear: %v", err)
			panic(err)
		}
		renderer.Copy(tex, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "renderer.Copy: %v", err)
			panic(err)
		}
		renderer.Present()
	}
}

// Causes Run() to return after processing any pending events.
func Quit() {
	sdl.PushEvent(&sdl.QuitEvent{Type: sdl.QUIT})
}
