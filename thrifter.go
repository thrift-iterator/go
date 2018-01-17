package thrifter

import (
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol"
	"errors"
	"io"
)

type Protocol int

var ProtocolBinary Protocol = 1

type Iterator interface {
	Error() error
	ReportError(operation string, err string)
	Reset(buf []byte)
	ReadMessageHeader() protocol.MessageHeader
	ReadMessage() protocol.Message
	ReadStructCB(func(fieldType protocol.TType, fieldId protocol.FieldId))
	ReadStructField() (fieldType protocol.TType, fieldId protocol.FieldId)
	ReadStruct() map[protocol.FieldId]interface{}
	SkipStruct() []byte
	ReadListHeader() (elemType protocol.TType, size int)
	ReadList() []interface{}
	SkipList() []byte
	ReadMapHeader() (keyType protocol.TType, elemType protocol.TType, size int)
	ReadMap() map[interface{}]interface{}
	SkipMap() []byte
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
	Read(ttype protocol.TType) interface{}
	ReaderOf(ttype protocol.TType) func() interface{}
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
	NewIterator(buf []byte) Iterator
	NewStream(buf []byte) Stream
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

func (cfg *frozenConfig) NewIterator(buf []byte) Iterator {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewIterator(buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) NewStream(buf []byte) Stream {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewStream(buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) Unmarshal(buf []byte, obj interface{}) error {
	msg, _ := obj.(*protocol.Message)
	if msg == nil {
		return errors.New("can only unmarshal protocol.Message")
	}
	iter := cfg.NewIterator(buf)
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
	stream := cfg.NewStream(nil)
	stream.WriteMessage(msg)
	if stream.Error() != nil {
		return nil, stream.Error()
	}
	return stream.Buffer(), nil
}

func (cfg *frozenConfig) NewDecoder(reader io.Reader) Decoder {
	if cfg.isFramed {
		switch cfg.protocol {
		case ProtocolBinary:
			return &framedDecoder{reader: reader, iter: cfg.NewIterator(nil)}
		}
	}
	panic("unsupported protocol")
}

var DefaultConfig = Config{Protocol: ProtocolBinary, IsFramed: true}.Froze()

func NewIterator(buf []byte) Iterator {
	return DefaultConfig.NewIterator(buf)
}

func NewStream(buf []byte) Stream {
	return DefaultConfig.NewStream(buf)
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
