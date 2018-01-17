package thrifter

import (
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol"
)

type Protocol int

var ProtocolBinary Protocol = 1

type Iterator interface {
	Error() error
	ReportError(operation string, err string)
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

type Config struct {
	Protocol Protocol
}

type API interface {
	NewIterator(buf []byte) Iterator
	NewStream(buf []byte) Stream
}

type frozenConfig struct {
	protocol Protocol
}

func (cfg Config) Froze() API {
	api := &frozenConfig{protocol: cfg.Protocol}
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

var DefaultConfig = Config{Protocol: ProtocolBinary}.Froze()

func NewIterator(buf []byte) Iterator {
	return DefaultConfig.NewIterator(buf)
}

func NewStream(buf []byte) Stream {
	return DefaultConfig.NewStream(buf)
}
