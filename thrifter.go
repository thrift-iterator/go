package thrifter

import (
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol"
	"errors"
	"io"
	"github.com/thrift-iterator/go/protocol/sbinary"
)

type Protocol int

var ProtocolBinary Protocol = 1

type Iterator interface {
	Error() error
	ReportError(operation string, err string)
	ReadMessageHeader() protocol.MessageHeader
	ReadMessage() protocol.Message
	SkipMessage(space []byte) []byte
	ReadStructCB(func(fieldType protocol.TType, fieldId protocol.FieldId))
	ReadStructField() (fieldType protocol.TType, fieldId protocol.FieldId)
	ReadStruct() map[protocol.FieldId]interface{}
	SkipStruct(space []byte) []byte
	ReadListHeader() (elemType protocol.TType, size int)
	ReadList() []interface{}
	SkipList(space []byte) []byte
	ReadMapHeader() (keyType protocol.TType, elemType protocol.TType, size int)
	ReadMap() map[interface{}]interface{}
	SkipMap(space []byte) []byte
	ReadBool() bool
	ReadInt8() int8
	ReadUInt8() uint8
	ReadInt16() int16
	ReadUInt16() uint16
	ReadInt32() int32
	ReadUInt32() uint32
	ReadInt64() int64
	ReadUInt64() uint64
	ReadFloat64() float64
	ReadString() string
	ReadBinary() []byte
	SkipBinary(space []byte) []byte
	Read(ttype protocol.TType) interface{}
	ReaderOf(ttype protocol.TType) func() interface{}
}

type BufferedIterator interface {
	Iterator
	Reset(buf []byte)
}

type Stream interface {
	Error() error
	ReportError(operation string, err string)
	Buffer() []byte
	WriteMessageHeader(header protocol.MessageHeader)
	WriteMessage(message protocol.Message)
	WriteListHeader(elemType protocol.TType, length int)
	WriteList(val []interface{})
	WriteStructField(fieldType protocol.TType, fieldId protocol.FieldId)
	WriteStructFieldStop()
	WriteStruct(val map[protocol.FieldId]interface{})
	WriteMapHeader(keyType protocol.TType, elemType protocol.TType, length int)
	WriteMap(val map[interface{}]interface{})
	WriterOf(sample interface{}) (protocol.TType, func(interface{}))
	WriteBool(val bool)
	WriteInt8(val int8)
	WriteUInt8(val uint8)
	WriteInt16(val int16)
	WriteUInt16(val uint16)
	WriteInt32(val int32)
	WriteUInt32(val uint32)
	WriteInt64(val int64)
	WriteUInt64(val uint64)
	WriteFloat64(val float64)
	WriteBinary(val []byte)
	WriteString(val string)
}

type Decoder interface {
	Decode(obj interface{}) error
}

type Config struct {
	Protocol Protocol
	IsFramed bool
}

type API interface {
	NewBufferedIterator(buf []byte) BufferedIterator
	NewBufferedStream(buf []byte) Stream
	NewIterator(reader io.Reader) Iterator
	Unmarshal(buf []byte, obj interface{}) error
	Marshal(obj interface{}) ([]byte, error)
	NewDecoder(reader io.Reader) Decoder
}

type frozenConfig struct {
	protocol Protocol
	isFramed bool
}

func (cfg Config) Froze() API {
	api := &frozenConfig{protocol: cfg.Protocol, isFramed: cfg.IsFramed}
	return api
}

func (cfg *frozenConfig) NewBufferedIterator(buf []byte) BufferedIterator {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewIterator(buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) NewBufferedStream(buf []byte) Stream {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewStream(buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) NewIterator(reader io.Reader) Iterator {
	switch cfg.protocol {
	case ProtocolBinary:
		return sbinary.NewIterator(reader)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) Unmarshal(buf []byte, obj interface{}) error {
	msg, _ := obj.(*protocol.Message)
	if msg == nil {
		return errors.New("can only unmarshal protocol.Message")
	}
	if cfg.isFramed {
		size := uint32(buf[3]) | uint32(buf[2])<<8 | uint32(buf[1])<<16 | uint32(buf[0])<<24
		buf = buf[4:4+size]
	}
	iter := cfg.NewBufferedIterator(buf)
	msgRead := iter.ReadMessage()
	if iter.Error() != nil {
		return iter.Error()
	}
	msg.Set(&msgRead)
	return nil
}

func (cfg *frozenConfig) Marshal(obj interface{}) ([]byte, error) {
	msg, isMsg := obj.(protocol.Message)
	if !isMsg {
		return nil, errors.New("can only unmarshal protocol.Message")
	}
	stream := cfg.NewBufferedStream(nil)
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
		return &framedDecoder{reader: reader, iter: cfg.NewBufferedIterator(nil)}
	} else {
		return &unframedDecoder{iter: cfg.NewIterator(reader)}
	}
}

var DefaultConfig = Config{Protocol: ProtocolBinary, IsFramed: true}.Froze()

func NewBufferedIterator(buf []byte) BufferedIterator {
	return DefaultConfig.NewBufferedIterator(buf)
}

func NewBufferedStream(buf []byte) Stream {
	return DefaultConfig.NewBufferedStream(buf)
}

func NewIterator(reader io.Reader) Iterator {
	return DefaultConfig.NewIterator(reader)
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
