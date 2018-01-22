package dynamic

import (
	"github.com/thrift-iterator/go/protocol"
	"unsafe"
	"github.com/thrift-iterator/go/spi"
)

type structEncoder struct {
	fields []structEncoderField
}

type structEncoderField struct {
	offset  uintptr
	fieldId protocol.FieldId
	encoder internalEncoder
}

func (encoder *structEncoder) encode(ptr unsafe.Pointer, stream spi.Stream) {
	stream.WriteStructHeader()
	for _, field := range encoder.fields {
		stream.WriteStructField(field.encoder.thriftType(), field.fieldId)
		field.encoder.encode(unsafe.Pointer(uintptr(ptr)+field.offset), stream)
	}
	stream.WriteStructFieldStop()
}

func (encoder *structEncoder) thriftType() protocol.TType {
	return protocol.TypeStruct
}