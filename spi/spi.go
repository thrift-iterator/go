package spi

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
)

type Iterator interface {
	Error() error
	Reset(reader io.Reader, buf []byte)
	ReportError(operation string, err string)
	ReadMessageHeader() protocol.MessageHeader
	ReadMessage() protocol.Message
	SkipMessage(space []byte) []byte
	ReadStructCB(func(fieldType protocol.TType, fieldId protocol.FieldId))
	ReadStructHeader()
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
	ReadInt() int
	ReadUint() uint
	ReadInt8() int8
	ReadUint8() uint8
	ReadInt16() int16
	ReadUint16() uint16
	ReadInt32() int32
	ReadUint32() uint32
	ReadInt64() int64
	ReadUint64() uint64
	ReadFloat64() float64
	ReadString() string
	ReadBinary() []byte
	SkipBinary(space []byte) []byte
	Discard(ttype protocol.TType)
	Read(ttype protocol.TType) interface{}
	ReaderOf(ttype protocol.TType) func() interface{}
}

type Stream interface {
	Error() error
	ReportError(operation string, err string)
	Reset(writer io.Writer)
	Flush()
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
	WriteInt(val int)
	WriteUint(val uint)
	WriteInt8(val int8)
	WriteUint8(val uint8)
	WriteInt16(val int16)
	WriteUint16(val uint16)
	WriteInt32(val int32)
	WriteUint32(val uint32)
	WriteInt64(val int64)
	WriteUint64(val uint64)
	WriteFloat64(val float64)
	WriteBinary(val []byte)
	WriteString(val string)
}

type ValEncoder interface {
	Encode(val interface{}, stream Stream)
}

type ValDecoder interface {
	Decode(val interface{}, iter Iterator)
}
