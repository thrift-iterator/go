package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/protocol/binary"
	"fmt"
	"github.com/thrift-iterator/go/spi"
)

type Iterator struct {
	spi.ValDecoderProvider
	real     *binary.Iterator
	reader   io.Reader
	// tmp buffer, because recorder is not always on, we can not use recorder to replace this one
	tmp      []byte
	// reusable recorder buffer
	space    []byte
	// recording read bytes into this buffer if not nil
	recorder []byte
	err      error
}

func NewIterator(provider spi.ValDecoderProvider, reader io.Reader, buf []byte) *Iterator {
	return &Iterator{
		ValDecoderProvider: provider,
		reader: reader, real: binary.NewIterator(provider, nil),
		tmp:    make([]byte, 10),
		space:  buf,
	}
}

func (iter *Iterator) Spawn() spi.Iterator {
	return binary.NewIterator(iter.ValDecoderProvider, nil)
}

func (iter *Iterator) Reset(reader io.Reader, buf []byte) {
	iter.reader = reader
	iter.err = nil
}

func (iter *Iterator) allocate(nBytes int) []byte {
	if len(iter.tmp) < nBytes {
		iter.tmp = make([]byte, nBytes)
	}
	return iter.tmp[:nBytes]
}

func (iter *Iterator) readSmall(nBytes int) []byte {
	tmp := iter.tmp[:nBytes]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.handleReadError(err, tmp)
		return tmp
	}
	if iter.recorder != nil {
		iter.recorder = append(iter.recorder, tmp...)
	}
	return tmp
}

func (iter *Iterator) readLarge(nBytes int) []byte {
	tmp := iter.allocate(nBytes)
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.handleReadError(err, tmp)
		return tmp
	}
	if iter.recorder != nil {
		iter.recorder = append(iter.recorder, tmp...)
	}
	return tmp
}

func (iter *Iterator) handleReadError(err error, tmp []byte) {
	for i := 0; i < len(tmp); i++ {
		tmp[i] = 0
	}
	iter.ReportError("read", err.Error())
}

func (iter *Iterator) Error() error {
	return iter.err
}

func (iter *Iterator) ReportError(operation string, err string) {
	if iter.err == nil {
		iter.err = fmt.Errorf("%s: %s", operation, err)
	}
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
	tmp := iter.readSmall(1)
	fieldType := protocol.TType(tmp[0])
	if fieldType == protocol.TypeStop {
		return protocol.TypeStop, 0
	}
	fieldId := protocol.FieldId(iter.ReadUint16())
	return fieldType, fieldId
}

func (iter *Iterator) ReadListHeader() (elemType protocol.TType, length int) {
	tmp := iter.readSmall(5)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadListHeader()
}

func (iter *Iterator) ReadMapHeader() (keyType protocol.TType, elemType protocol.TType, length int) {
	tmp := iter.readSmall(6)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadMapHeader()
}

func (iter *Iterator) ReadBool() bool {
	tmp := iter.readSmall(1)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadBool()
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUint8())
}

func (iter *Iterator) ReadUint8() uint8 {
	tmp := iter.readSmall(1)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint8()
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadUint16())
}

func (iter *Iterator) ReadUint16() uint16 {
	tmp := iter.readSmall(2)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint16()
}

func (iter *Iterator) ReadInt32() int32 {
	return int32(iter.ReadUint32())
}

func (iter *Iterator) ReadUint32() uint32 {
	tmp := iter.readSmall(4)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint32()
}

func (iter *Iterator) ReadInt64() int64 {
	tmp := iter.readSmall(8)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadInt64()
}

func (iter *Iterator) ReadUint64() uint64 {
	tmp := iter.readSmall(8)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint64()
}

func (iter *Iterator) ReadInt() int {
	return int(iter.ReadInt64())
}

func (iter *Iterator) ReadUint() uint {
	return uint(iter.ReadUint64())
}

func (iter *Iterator) ReadFloat64() float64 {
	tmp := iter.readSmall(8)
	iter.real.Reset(nil, tmp)
	return iter.real.ReadFloat64()
}

func (iter *Iterator) ReadString() string {
	size := iter.ReadUint32()
	tmp := iter.readLarge(int(size))
	return string(tmp)
}

func (iter *Iterator) ReadBinary() []byte {
	size := iter.ReadUint32()
	tmp := make([]byte, size)
	copy(tmp, iter.readLarge(int(size)))
	return tmp
}

func getTypeSize(elemType protocol.TType) int {
	switch elemType {
	case protocol.TypeBool, protocol.TypeI08:
		return 1
	case protocol.TypeI16:
		return 2
	case protocol.TypeI32:
		return 4
	case protocol.TypeI64, protocol.TypeDouble:
		return 8
	}
	return 0
}
