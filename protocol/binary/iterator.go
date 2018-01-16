package binary

import (
	"fmt"
	"github.com/thrift-iterator/go/protocol"
	"math"
)

type Iterator struct {
	buf   []byte
	Error error
}

var typeSizes = map[protocol.TType]int{
	protocol.BOOL:   1,
	protocol.I08:    1,
	protocol.I16:    2,
	protocol.I32:    4,
	protocol.I64:    8,
	protocol.DOUBLE: 8,
}

func NewIterator(buf []byte) *Iterator {
	return &Iterator{buf: buf}
}

func (iter *Iterator) ReportError(operation string, err string) {
	if iter.Error == nil {
		iter.Error = fmt.Errorf("%s: %s", operation, err)
	}
}

func (iter *Iterator) ReadMessageHeader() protocol.MessageHeader {
	versionAndMessageType := iter.ReadInt32()
	messageType := protocol.TMessageType(versionAndMessageType & 0x0ff)
	version := protocol.Version(int64(int64(versionAndMessageType) & 0xffff0000))
	messageName := iter.ReadString()
	seqId := protocol.SeqId(iter.ReadInt32())
	return protocol.MessageHeader{
		Version:     version,
		MessageName: messageName,
		MessageType: messageType,
		SeqId:       seqId,
	}
}

func (iter *Iterator) ReadStructCB(cb func(fieldType protocol.TType, fieldId protocol.FieldId)) {
	for iter.buf[0] != 0 {
		fieldType := iter.buf[0]
		fieldId := uint16(iter.buf[2]) | uint16(iter.buf[1])<<8
		iter.buf = iter.buf[3:]
		cb(protocol.TType(fieldType), protocol.FieldId(fieldId))
	}
	iter.buf = iter.buf[1:]
}

func (iter *Iterator) ReadStruct() (protocol.TType, protocol.FieldId) {
	fieldType := iter.buf[0]
	if fieldType == 0 {
		iter.buf = iter.buf[1:]
		return protocol.TType(fieldType), 0
	}
	fieldId := uint16(iter.buf[2]) | uint16(iter.buf[1])<<8
	iter.buf = iter.buf[3:]
	return protocol.TType(fieldType), protocol.FieldId(fieldId)
}

func (iter *Iterator) SkipStruct() []byte {
	bufBeforeSkip := iter.buf
	skippedBytes := 0
	for {
		fieldType := protocol.TType(iter.buf[0])
		if fieldType == 0 {
			iter.buf = iter.buf[1:]
			skippedBytes += 1
			return bufBeforeSkip[:skippedBytes]
		}
		switch fieldType {
		case protocol.BOOL, protocol.I08:
			iter.buf = iter.buf[4:]
			skippedBytes += 4
		case protocol.I16:
			iter.buf = iter.buf[5:]
			skippedBytes += 5
		case protocol.I32:
			iter.buf = iter.buf[7:]
			skippedBytes += 7
		case protocol.I64, protocol.DOUBLE:
			iter.buf = iter.buf[11:]
			skippedBytes += 11
		default:
			panic("unsupported type")
		}
	}
}

func (iter *Iterator) ReadList() (protocol.TType, int) {
	b := iter.buf
	elemType := b[0]
	length := uint32(b[4]) | uint32(b[3])<<8 | uint32(b[2])<<16 | uint32(b[1])<<24
	iter.buf = iter.buf[5:]
	return protocol.TType(elemType), int(length)
}

func (iter *Iterator) SkipList() []byte {
	bufBeforeSkip := iter.buf
	elemType := protocol.TType(bufBeforeSkip[0])
	length := uint32(bufBeforeSkip[4]) | uint32(bufBeforeSkip[3])<<8 | uint32(bufBeforeSkip[2])<<16 | uint32(bufBeforeSkip[1])<<24
	switch elemType {
	case protocol.BOOL, protocol.I08:
		size := 5 + length
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.I16:
		size := 5 + length*2
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.I32:
		size := 5 + length*4
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.I64, protocol.DOUBLE:
		size := 5 + length*8
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.STRING:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.ReadBinary())
			skippedBytes += 4
		}
		iter.buf = bufBeforeSkip[skippedBytes:]
		return bufBeforeSkip[:skippedBytes]
	case protocol.LIST:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipList())
		}
		iter.buf = bufBeforeSkip[skippedBytes:]
		return bufBeforeSkip[:skippedBytes]
	case protocol.MAP:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipMap())
		}
		iter.buf = bufBeforeSkip[skippedBytes:]
		return bufBeforeSkip[:skippedBytes]
	case protocol.STRUCT:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipStruct())
		}
		iter.buf = bufBeforeSkip[skippedBytes:]
		return bufBeforeSkip[:skippedBytes]
	}
	panic("unsupported type")
}

func (iter *Iterator) ReadMap() (protocol.TType, protocol.TType, int) {
	b := iter.buf
	keyType := b[0]
	elemType := b[1]
	length := uint32(b[5]) | uint32(b[4])<<8 | uint32(b[3])<<16 | uint32(b[2])<<24
	iter.buf = iter.buf[6:]
	return protocol.TType(keyType), protocol.TType(elemType), int(length)
}

func (iter *Iterator) SkipMap() []byte {
	b := iter.buf
	keyType := protocol.TType(b[0])
	elemType := protocol.TType(b[1])
	length := uint32(b[5]) | uint32(b[4])<<8 | uint32(b[3])<<16 | uint32(b[2])<<24
	keySize := getTypeSize(keyType)
	elemSize := getTypeSize(elemType)
	if keySize != 0 && elemSize != 0 {
		size := 6 + int(length)*(elemSize+keySize)
		skipped := b[:size]
		iter.buf = b[size:]
		return skipped
	}
	panic("unsupported type")
}

func getTypeSize(elemType protocol.TType) int {
	switch elemType {
	case protocol.BOOL, protocol.I08:
		return 1
	case protocol.I16:
		return 2
	case protocol.I32:
		return 4
	case protocol.I64, protocol.DOUBLE:
		return 8
	}
	return 0
}

func (iter *Iterator) ReadBool() bool {
	return iter.ReadUInt8() == 1
}

func (iter *Iterator) ReadUInt8() uint8 {
	b := iter.buf
	value := b[0]
	iter.buf = iter.buf[1:]
	return value
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUInt8())
}

func (iter *Iterator) ReadUInt16() uint16 {
	b := iter.buf
	value := uint16(b[1]) | uint16(b[0])<<8
	iter.buf = iter.buf[2:]
	return value
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadUInt16())
}

func (iter *Iterator) ReadUInt32() uint32 {
	b := iter.buf
	value := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	iter.buf = iter.buf[4:]
	return value
}

func (iter *Iterator) ReadInt32() int32 {
	return int32(iter.ReadUInt32())
}

func (iter *Iterator) ReadInt64() int64 {
	return int64(iter.ReadUInt64())
}

func (iter *Iterator) ReadUInt64() uint64 {
	b := iter.buf
	value := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	iter.buf = iter.buf[8:]
	return value
}

func (iter *Iterator) ReadFloat64() float64 {
	return math.Float64frombits(iter.ReadUInt64())
}

func (iter *Iterator) ReadString() string {
	length := iter.ReadUInt32()
	value := string(iter.buf[:length])
	iter.buf = iter.buf[length:]
	return value
}

func (iter *Iterator) ReadBinary() []byte {
	length := iter.ReadUInt32()
	value := iter.buf[:length]
	iter.buf = iter.buf[length:]
	return value
}
