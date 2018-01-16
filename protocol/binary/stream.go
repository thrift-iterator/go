package binary

import (
	"math"
)

type Stream struct {
	buf   []byte
	Error error
}

func NewStream(buf []byte) *Stream {
	return &Stream{
		buf: buf,
	}
}

func (stream *Stream) Buffer() []byte {
	return stream.buf
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