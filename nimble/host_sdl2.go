// Plaform-dependent routines of nimble.
// This file implements them using SDL 2.

package nimble

// typedef unsigned char Uint8;
// void getSoundSamplesAdaptor(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"os"
	"reflect"
	"runtime"
	"unsafe"
)

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

// win==nil requests a full-screen window
func Run(win WindowSpec) int {
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
	var window *sdl.Window
	if win != nil {
		// Partial screen
		winWidth, winHeight := win.Size()
		window, err = sdl.CreateWindow(win.Title(), sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			int(winWidth), int(winHeight), sdl.WINDOW_SHOWN)
	} else {
		// Full screen
		window, err = sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			0, 0, sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN_DESKTOP)
	}
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

	// Loop until quit
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return 0
			case *sdl.MouseMotionEvent:
				// Go equivalent of SDL_PRESSED seems to be missing, so compare with zero.
				if e.State != 0 {
					forwardMouseEvent(MouseDrag, int32(e.X), int32(e.Y))
				} else {
					forwardMouseEvent(MouseMove, int32(e.X), int32(e.Y))
				}
			case *sdl.MouseButtonEvent:
				switch e.Type {
				case sdl.MOUSEBUTTONDOWN:
					forwardMouseEvent(MouseDown, int32(e.X), int32(e.Y))
				case sdl.MOUSEBUTTONUP:
					forwardMouseEvent(MouseUp, int32(e.X), int32(e.Y))
				}
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
					forwardKeyEvent(k)
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

// Quit causes Run() to return after processing any pending events.
func Quit() {
	sdl.PushEvent(&sdl.QuitEvent{Type: sdl.QUIT})
}

// ShowCursor causes the cursor to be shown or hidden.
func ShowCursor(show bool) {
	if show {
		sdl.ShowCursor(1)
	} else {
		sdl.ShowCursor(0)
	}
}

// A Font
type Font ttf.Font

// OpenFont opens creates a font object and returns a pointer to it.
func OpenFont(filename string, size int) (*Font, error) {
	if !ttf.WasInit() {
		ttf.Init()
	}
	f, err := ttf.OpenFont(filename, size)
	return (*Font)(f), err
}

// Close frees resources used by the given Font object.
func (f Font) Close() {
	f.Close()
}

// DrawText draws the given text at (x,y) on pm, in the given color and font.
// The return value indicates the width and height of a bounding box for the text.
func (pm *PixMap) DrawText(x, y int32, text string, color Pixel, font *Font) (width, height int32) {
	if text == "" {
		return 0, font.Height()
	}
	c := sdl.Color{R: uint8(color >> redShift),
		G: uint8(color >> greenShift),
		B: uint8(color >> blueShift),
		A: uint8(color >> alphaShift)}
	tmp, err := (*ttf.Font)(font).RenderUTF8_Solid(text, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "RenderUTF8_Solid: %v", err)
		panic(err)
	}
	defer tmp.Free()
	width, height = tmp.W, tmp.H
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&pm.buf))
	dst, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(sliceHeader.Data),
		int(pm.width), int(pm.height), 32, int(pm.vstride*4),
		0xFF<<redShift, 0xFF<<greenShift, 0xFF<<blueShift, 0xFF<<alphaShift)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateRGBSurfaceFrom: %v", err)
		panic(err)
	}
	defer dst.Free()
	err = tmp.Blit(nil, dst, &sdl.Rect{X: x, Y: y, W: tmp.W, H: tmp.H})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Blit: %v\n", err)
		panic(err)
	}
	return
}

// Height returns the nominal height of the font.
func (f *Font) Height() int32 {
	return int32((*ttf.Font)(f).Height())
}

// Size returns the width and height of a bounding box for the text.
func (f *Font) Size(text string) (width, height int32) {
	w, h, err := (*ttf.Font)(f).SizeUTF8(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SizeUTF8: %v\n", err)
		panic(err)
	}
	return int32(w), int32(h)
}
