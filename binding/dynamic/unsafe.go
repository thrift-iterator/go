package dynamic

import (
	"unsafe"
	"github.com/thrift-iterator/go/spi"
)

type internalDecoder interface {
	decode(ptr unsafe.Pointer, iter spi.Iterator)
}

type valDecoderAdapter struct {
	decoder internalDecoder
}

func (decoder *valDecoderAdapter) Decode(val interface{}, iter spi.Iterator) {
	ptr := (*emptyInterface)(unsafe.Pointer(&val)).word
	decoder.decoder.decode(ptr, iter)
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}