# NimbleDraw
Go package for simple 2D video games.

Package NimbleDraw provides a lightweight graphics interface for my video games.
The games typically involve unusual rendering not handled by hardware accelerators.  
Hence the interface focuses on providing simple efficient access to bitmaps in memory.

The package is also intended to be a thin wrapper around a host platform.
The initial port is against [go-sdl2](https://github.com/veandco/go-sdl2)
