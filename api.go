package thrifter

import (
	"io"
)

type Protocol int

var ProtocolBinary Protocol = 1
var ProtocolCompact Protocol = 2

type Decoder interface {
	Decode(obj interface{}) error
}

type Encoder interface {
	Encode(obj interface{}) error
}

type Config struct {
	Protocol         Protocol
	IsFramed         bool
}

type API interface {
	// NewStream is low level streaming api
	NewStream(writer io.Writer, buf []byte) Stream
	// NewIterator is low level streaming api
	NewIterator(reader io.Reader, buf []byte) Iterator
	Unmarshal(buf []byte, obj interface{}) error
	Marshal(obj interface{}) ([]byte, error)
	NewDecoder(reader io.Reader) Decoder
	NewEncoder(writer io.Writer) Encoder
}

var DefaultConfig = Config{Protocol: ProtocolBinary, IsFramed: true}.Froze()

func NewStream(writer io.Writer, buf []byte) Stream {
	return DefaultConfig.NewStream(writer, buf)
}

func NewIterator(reader io.Reader, buf []byte) Iterator {
	return DefaultConfig.NewIterator(reader, buf)
}

func Unmarshal(buf []byte, obj interface{}) error {
	return DefaultConfig.Unmarshal(buf, obj)
}

func Marshal(obj interface{}) ([]byte, error) {
	return DefaultConfig.Marshal(obj)
}

func NewDecoder(reader io.Reader) Decoder {
	return DefaultConfig.NewDecoder(reader)
}

func NewEncoder(writer io.Writer) Encoder {
	return DefaultConfig.NewEncoder(writer)
}
