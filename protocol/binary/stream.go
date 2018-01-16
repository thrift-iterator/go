package binary

type Stream struct {
	buf []byte
	Error error
}

func NewStream(buf []byte) *Stream {
	return &Stream {
		buf: buf,
	}
}

func (stream *Stream) WriteInt8(val int8) {
	stream.buf = append(stream.buf, byte(val))
}

func (stream *Stream) Buffer() []byte {
	return stream.buf
}