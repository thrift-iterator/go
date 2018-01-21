package dynamic

import (
	"unsafe"
	"github.com/thrift-iterator/go/spi"
)

type binaryEncoder struct {
}

func (encoder *binaryEncoder) encode(ptr unsafe.Pointer, iter spi.Stream) {
	iter.WriteBinary(*(*[]byte)(ptr))
}

type stringEncoder struct {
}

func (encoder *stringEncoder) encode(ptr unsafe.Pointer, iter spi.Stream) {
	iter.WriteString(*(*string)(ptr))
}

type boolEncoder struct {
}

func (encoder *boolEncoder) encode(ptr unsafe.Pointer, iter spi.Stream) {
	iter.WriteBool(*(*bool)(ptr))
}

type float64Encoder struct {
}

func (encoder *float64Encoder) encode(ptr unsafe.Pointer, iter spi.Stream) {
	iter.WriteFloat64(*(*float64)(ptr))
}

type int32Encoder struct {
}

func (encoder *int32Encoder) encode(ptr unsafe.Pointer, iter spi.Stream) {
	iter.WriteInt32(*(*int32)(ptr))
}
