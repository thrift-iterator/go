package binary

import (
	"fmt"
	"github.com/thrift-iterator/go/protocol"
	"math"
	"io"
	"github.com/thrift-iterator/go/spi"
)

type Iterator struct {
	spi.ValDecoderProvider
	buf []byte
	err error
}

func NewIterator(provider spi.ValDecoderProvider, buf []byte) *Iterator {
	return &Iterator{ValDecoderProvider: provider, buf: buf}
}

func (iter *Iterator) Error() error {
	return iter.err
}

func (iter *Iterator) ReportError(operation string, err string) {
	if iter.err == nil {
		iter.err = fmt.Errorf("%s: %s", operation, err)
	}
}

func (iter *Iterator) Reset(reader io.Reader, buf []byte) {
	iter.buf = buf
	iter.err = nil
}

const version1 = 0x80010000

func (iter *Iterator) ReadMessageHeader() protocol.MessageHeader {
	versionAndMessageType := iter.ReadInt32()
	messageType := protocol.TMessageType(versionAndMessageType & 0x0ff)
	version := int64(int64(versionAndMessageType) & 0xffff0000)
	if version != version1 {
		iter.ReportError("ReadMessageHeader", "unexpected version")
		return protocol.MessageHeader{}
	}
	messageName := iter.ReadString()
	seqId := protocol.SeqId(iter.ReadInt32())
	return protocol.MessageHeader{
		MessageName: messageName,
		MessageType: messageType,
		SeqId:       seqId,
	}
}

func (iter *Iterator) ReadStructHeader() {
	// noop
}

func (iter *Iterator) ReadStructField() (protocol.TType, protocol.FieldId) {
	fieldType := iter.buf[0]
	if fieldType == 0 {
		iter.buf = iter.buf[1:]
		return protocol.TType(fieldType), 0
	}
	fieldId := uint16(iter.buf[2]) | uint16(iter.buf[1])<<8
	iter.buf = iter.buf[3:]
	return protocol.TType(fieldType), protocol.FieldId(fieldId)
}

func (iter *Iterator) ReadListHeader() (protocol.TType, int) {
	b := iter.buf
	elemType := b[0]
	length := uint32(b[4]) | uint32(b[3])<<8 | uint32(b[2])<<16 | uint32(b[1])<<24
	iter.buf = iter.buf[5:]
	return protocol.TType(elemType), int(length)
}

func (iter *Iterator) ReadMapHeader() (protocol.TType, protocol.TType, int) {
	b := iter.buf
	keyType := b[0]
	elemType := b[1]
	length := uint32(b[5]) | uint32(b[4])<<8 | uint32(b[3])<<16 | uint32(b[2])<<24
	iter.buf = iter.buf[6:]
	return protocol.TType(keyType), protocol.TType(elemType), int(length)
}

func (iter *Iterator) ReadBool() bool {
	return iter.ReadUint8() == 1
}

func (iter *Iterator) ReadUint8() uint8 {
	b := iter.buf
	value := b[0]
	iter.buf = iter.buf[1:]
	return value
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUint8())
}

func (iter *Iterator) ReadUint16() uint16 {
	b := iter.buf
	value := uint16(b[1]) | uint16(b[0])<<8
	iter.buf = iter.buf[2:]
	return value
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadUint16())
}

func (iter *Iterator) ReadUint32() uint32 {
	b := iter.buf
	value := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	iter.buf = iter.buf[4:]
	return value
}

func (iter *Iterator) ReadInt32() int32 {
	return int32(iter.ReadUint32())
}

func (iter *Iterator) ReadInt64() int64 {
	return int64(iter.ReadUint64())
}

func (iter *Iterator) ReadUint64() uint64 {
	b := iter.buf
	value := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	iter.buf = iter.buf[8:]
	return value
}

func (iter *Iterator) ReadInt() int {
	return int(iter.ReadInt64())
}

func (iter *Iterator) ReadUint() uint {
	return uint(iter.ReadUint64())
}

func (iter *Iterator) ReadFloat64() float64 {
	return math.Float64frombits(iter.ReadUint64())
}

func (iter *Iterator) ReadString() string {
	length := iter.ReadUint32()
	value := string(iter.buf[:length])
	iter.buf = iter.buf[length:]
	return value
}

func (iter *Iterator) ReadBinary() []byte {
	length := iter.ReadUint32()
	value := iter.buf[:length]
	iter.buf = iter.buf[length:]
	return value
}