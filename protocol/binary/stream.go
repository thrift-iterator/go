package binary

import (
	"math"
	"github.com/thrift-iterator/go/protocol"
	"fmt"
)

type Stream struct {
	buf   []byte
	err error
}

func NewStream(buf []byte) *Stream {
	return &Stream{
		buf: buf,
	}
}

func (stream *Stream) Error() error {
	return stream.err
}

func (stream *Stream) ReportError(operation string, err string) {
	if stream.err == nil {
		stream.err = fmt.Errorf("%s: %s", operation, err)
	}
}

func (stream *Stream) Buffer() []byte {
	return stream.buf
}

func (stream *Stream) WriteListHeader(elemType protocol.TType, length int) {
	stream.buf = append(stream.buf, byte(elemType),
		byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
}

func (stream *Stream) WriteList(val []interface{}) {
	if len(val) == 0 {
		stream.ReportError("WriteList", "input is empty slice, can not tell element type")
		return
	}
	switch val[0].(type) {
	case int64:
		stream.WriteListHeader(protocol.I64, len(val))
		for _, elem := range val {
			stream.WriteInt64(elem.(int64))
		}
	default:
		panic("unsupported type")
	}
}

func (stream *Stream) WriteStructField(fieldType protocol.TType, fieldId protocol.FieldId) {
	stream.buf = append(stream.buf, byte(fieldType), byte(fieldId>>8), byte(fieldId))
}

func (stream *Stream) WriteStructFieldStop() {
	stream.buf = append(stream.buf, byte(protocol.STOP))
}

func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.WriteUInt8(1)
	} else {
		stream.WriteUInt8(0)
	}
}

func (stream *Stream) WriteInt8(val int8) {
	stream.WriteUInt8(uint8(val))
}

func (stream *Stream) WriteUInt8(val uint8) {
	stream.buf = append(stream.buf, byte(val))
}

func (stream *Stream) WriteInt16(val int16) {
	stream.WriteUInt16(uint16(val))
}

func (stream *Stream) WriteUInt16(val uint16) {
	stream.buf = append(stream.buf, byte(val>>8), byte(val))
}

func (stream *Stream) WriteInt32(val int32) {
	stream.WriteUInt32(uint32(val))
}

func (stream *Stream) WriteUInt32(val uint32) {
	stream.buf = append(stream.buf, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

func (stream *Stream) WriteInt64(val int64) {
	stream.WriteUInt64(uint64(val))
}

func (stream *Stream) WriteUInt64(val uint64) {
	stream.buf = append(stream.buf,
		byte(val>>56), byte(val>>48), byte(val>>40), byte(val>>32),
		byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

func (stream *Stream) WriteFloat64(val float64) {
	stream.WriteUInt64(math.Float64bits(val))
}

func (stream *Stream) WriteBinary(val []byte) {
	stream.WriteUInt32(uint32(len(val)))
	stream.buf = append(stream.buf, val...)
}

func (stream *Stream) WriteString(val string) {
	stream.WriteUInt32(uint32(len(val)))
	stream.buf = append(stream.buf, val...)
}