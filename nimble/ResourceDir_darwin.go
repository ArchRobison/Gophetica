// +build bundle,darwin

package nimble 

// #include <mach-o/dyld.h>
import "C"

import "path"

var resourceDir string

func init() {
    bufSize := C.uint32_t(0)
    C._NSGetExecutablePath(nil,&bufSize)
    buf := make([]C.char, bufSize)
    C._NSGetExecutablePath(&buf[0],&bufSize)
    resourceDir = path.Dir(C.GoString(&buf[0]))+"/../Resources/"
}
