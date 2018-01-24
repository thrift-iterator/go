package reflection

import (
	"unsafe"
	"github.com/thrift-iterator/go/spi"
	"reflect"
	"github.com/thrift-iterator/go/protocol"
)

type sliceEncoder struct {
	sliceType   reflect.Type
	elemType    reflect.Type
	elemEncoder internalEncoder
}

func (encoder *sliceEncoder) encode(ptr unsafe.Pointer, stream spi.Stream) {
	slice := (*sliceHeader)(ptr)
	stream.WriteListHeader(encoder.elemEncoder.thriftType(), slice.Len)
	offset := uintptr(slice.Data)
	for i := 0; i < slice.Len; i++ {
		encoder.elemEncoder.encode(unsafe.Pointer(offset), stream)
		offset += encoder.elemType.Size()
	}
}

func (encoder *sliceEncoder) thriftType() protocol.TType {
	return protocol.TypeList
}