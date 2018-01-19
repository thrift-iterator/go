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

type int8Decoder struct {
}

func (decoder *int8Decoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*int8)(ptr) = iter.ReadInt8()
}

type uint8Decoder struct {
}

func (decoder *uint8Decoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*uint8)(ptr) = iter.ReadUint8()
}

type int16Decoder struct {
}

func (decoder *int16Decoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*int16)(ptr) = iter.ReadInt16()
}

type uint16Decoder struct {
}

func (decoder *uint16Decoder) decode(ptr unsafe.Pointer, iter spi.Iterator) {
	*(*uint16)(ptr) = iter.ReadUint16()
}

