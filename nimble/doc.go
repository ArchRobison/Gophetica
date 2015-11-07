/*
nimble is a Go package for simple 2D video games.

It's a Go incarnation of a lightweight graphics interface that I use 
for video games.  The games typically involve unusual rendering 
not handled by hardware accelerators.  Hence the interface focuses 
on providing simple efficient access to bitmaps in memory.

The other purpose of the package is be a thin wrapper that can be wrapped 
around a variety of host platforms.  The initial port is against 
go-sdl2 (https://github.com/veandco/go-sdl2).

When used for production, it should be built with `-tags=release`.
*/
package nimble
