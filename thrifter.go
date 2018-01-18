package thrifter

import (
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol"
	"errors"
	"io"
	"github.com/thrift-iterator/go/protocol/sbinary"
	"github.com/thrift-iterator/go/protocol/compact"
	"reflect"
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
	Protocol Protocol
	IsFramed bool
	Decoders map[reflect.Type]ValDecoder
	Encoders map[reflect.Type]ValEncoder
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

type frozenConfig struct {
	protocol Protocol
	isFramed bool
	encoders map[reflect.Type]ValEncoder
	decoders map[reflect.Type]ValDecoder
}

func (cfg Config) Froze() API {
	decoders := cfg.Decoders
	if decoders == nil {
		decoders = map[reflect.Type]ValDecoder{}
	}
	encoders := cfg.Encoders
	if encoders == nil {
		encoders = map[reflect.Type]ValEncoder{}
	}
	api := &frozenConfig{
		protocol: cfg.Protocol,
		isFramed: cfg.IsFramed,
		encoders: encoders,
		decoders: decoders,
	}
	return api
}

func (cfg *frozenConfig) NewStream(writer io.Writer, buf []byte) Stream {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewStream(writer, buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) NewIterator(reader io.Reader, buf []byte) Iterator {
	switch cfg.protocol {
	case ProtocolBinary:
		if reader != nil {
			return sbinary.NewIterator(reader, buf)
		}
		return binary.NewIterator(buf)
	case ProtocolCompact:
		return compact.NewIterator(buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) Unmarshal(buf []byte, obj interface{}) error {
	decoder := cfg.decoders[reflect.TypeOf(obj)]
	if decoder == nil {
		decoder = msgDecoderInstance
	}
	if cfg.isFramed {
		size := uint32(buf[3]) | uint32(buf[2])<<8 | uint32(buf[1])<<16 | uint32(buf[0])<<24
		buf = buf[4:4+size]
	}
	iter := cfg.NewIterator(nil, buf)
	decoder.Decode(obj, iter)
	if iter.Error() != nil {
		return iter.Error()
	}
	return nil
}

func (cfg *frozenConfig) Marshal(obj interface{}) ([]byte, error) {
	msg, isMsg := obj.(protocol.Message)
	if !isMsg {
		return nil, errors.New("can only unmarshal protocol.Message")
	}
	stream := cfg.NewStream(nil, nil)
	stream.WriteMessage(msg)
	if stream.Error() != nil {
		return nil, stream.Error()
	}
	buf := stream.Buffer()
	if cfg.isFramed {
		size := len(buf)
		buf = append([]byte{
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		}, buf...)
	}
	return buf, nil
}

func (cfg *frozenConfig) NewDecoder(reader io.Reader) Decoder {
	if cfg.isFramed {
		return &framedDecoder{reader: reader, iter: cfg.NewIterator(nil, nil)}
	} else {
		return &unframedDecoder{iter: cfg.NewIterator(reader, make([]byte, 256))}
	}
}

func (cfg *frozenConfig) NewEncoder(writer io.Writer) Encoder {
	if cfg.isFramed {
		return &framedEncoder{writer: writer, stream: cfg.NewStream(nil, nil)}
	} else {
		return &unframedEncoder{stream: cfg.NewStream(writer, nil)}
	}
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
