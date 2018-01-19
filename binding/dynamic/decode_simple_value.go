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

type boolDecoder struct {
}

func (decoder *boolDecoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*bool)(ptr) = iter.ReadBool()
}

type float64Decoder struct {
}

func (decoder *float64Decoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*float64)(ptr) = iter.ReadFloat64()
}
