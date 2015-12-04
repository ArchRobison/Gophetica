// +build bundle,darwin

package nimble

// #include <mach-o/dyld.h>
// #include <CoreServices/CoreServices.h>
//
// void getApplicationSupportFolder(char path[PATH_MAX]) {
//     FSRef ref;
//     OSType folderType = kApplicationSupportFolderType;
//     FSFindFolder( kUserDomain, (OSType)kApplicationSupportFolderType, kCreateFolder, &ref );
//     FSRefMakePath(&ref, (UInt8*)path, PATH_MAX);
// }
// #cgo LDFLAGS: -framework CoreServices
import "C"

import "path"

var (
	resourceDir = getResourceDir()
	recordDir   = getRecordDir()
)

func getResourceDir() string {
	bufSize := C.uint32_t(0)
	C._NSGetExecutablePath(nil, &bufSize)
	buf := make([]C.char, bufSize)
	C._NSGetExecutablePath(&buf[0], &bufSize)
	return path.Dir(C.GoString(&buf[0])) + "/../Resources/"
}

func getRecordDir() string {
	buf := make([]C.char, C.PATH_MAX)
	C.getApplicationSupportFolder(&buf[0])
	return C.GoString(&buf[0]) + "/"
}

func init() {
	resourceDir = getResourceDir()
	recordDir = getRecordDir()
}
