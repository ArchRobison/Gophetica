// +build !bundle

package nimble

import (
	"testing"
)

func TestReadPixMap(t *testing.T) {
	pm, err := ReadPixMap("ReadPixMap_test.png")
	assert(err == nil, "ReadPixMap/err", t)
	assert(pm.Width() == 5, "ReadPixMap/Width", t)
	assert(pm.Height() == 3, "ReadPixMap/Height", t)
	assert(pm.Pixel(0, 0) == RGB(0, 0, 0), "ReadPixMap(0,0)", t)
	assert(pm.Pixel(3, 0) == RGB(0, 1, 0), "ReadPixMap(3,0)", t)
	assert(pm.Pixel(0, 2) == RGB(1, 0, 0), "ReadPixMap(0,2)", t)
	assert(pm.Pixel(2, 2) == RGB(1, 1, 0), "ReadPixMap(2,2)", t)
	assert(pm.Pixel(4, 2) == RGB(0, 0, 1), "ReadPixMap(4,2)", t)
	assert(pm.Pixel(4, 0) == Gray(128./255.), "ReadPixMap(4,0)", t)
}
