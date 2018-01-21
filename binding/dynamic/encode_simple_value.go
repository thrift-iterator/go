package dynamic

import (
	"unsafe"
	"github.com/thrift-iterator/go/spi"
)

type boolEncoder struct {
}

func (encoder *boolEncoder) encode(ptr unsafe.Pointer, iter spi.Stream) {
	iter.WriteBool(*(*bool)(ptr))
}
