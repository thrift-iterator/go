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
	SkipMessageHeader(space []byte) []byte
	ReadStructHeader()
	ReadStructField() (fieldType protocol.TType, fieldId protocol.FieldId)
	SkipStruct(space []byte) []byte
	ReadListHeader() (elemType protocol.TType, size int)
	SkipList(space []byte) []byte
	ReadMapHeader() (keyType protocol.TType, elemType protocol.TType, size int)
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
	Skip(ttype protocol.TType, space []byte) []byte
	Discard(ttype protocol.TType)
}

type Stream interface {
	Error() error
	ReportError(operation string, err string)
	Reset(writer io.Writer)
	Flush()
	Buffer() []byte
	Write(buf []byte) error
	WriteMessageHeader(header protocol.MessageHeader)
	WriteListHeader(elemType protocol.TType, length int)
	WriteStructHeader()
	WriteStructField(fieldType protocol.TType, fieldId protocol.FieldId)
	WriteStructFieldStop()
	WriteMapHeader(keyType protocol.TType, elemType protocol.TType, length int)
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
