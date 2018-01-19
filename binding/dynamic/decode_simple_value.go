package dynamic

import (
	"github.com/thrift-iterator/go/spi"
	"unsafe"
)

type binaryDecoder struct {
}

func (decoder *binaryDecoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*[]byte)(ptr) = iter.ReadBinary()
}
