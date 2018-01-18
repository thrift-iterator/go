package compact

import (
	"fmt"
	"github.com/thrift-iterator/go/protocol"
	"math"
	"io"
	"encoding/binary"
)

type Iterator struct {
	buf              []byte
	err              error
	fieldIdStack     []protocol.FieldId
	lastFieldId      protocol.FieldId
	consumed         int
	pendingBoolField uint8
}

func NewIterator(buf []byte) *Iterator {
	return &Iterator{buf: buf}
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

func (iter *Iterator) consume(nBytes int) {
	iter.buf = iter.buf[nBytes:]
	iter.consumed += nBytes
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

func (iter *Iterator) ReadMessage() protocol.Message {
	header := iter.ReadMessageHeader()
	return protocol.Message{
		MessageHeader: header,
		Arguments:     iter.ReadStruct(),
	}
}

func (iter *Iterator) ReadStructCB(cb func(fieldType protocol.TType, fieldId protocol.FieldId)) {
	iter.ReadStructHeader()
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.STOP {
			return
		}
		cb(fieldType, fieldId)
	}
}

func (iter *Iterator) ReadStructHeader() {
	iter.fieldIdStack = append(iter.fieldIdStack, iter.lastFieldId)
	iter.lastFieldId = 0
}

func (iter *Iterator) ReadStructField() (protocol.TType, protocol.FieldId) {
	firstByte := iter.buf[0]
	iter.consume(1)
	if firstByte == 0 {
		iter.lastFieldId = iter.fieldIdStack[len(iter.fieldIdStack)-1]
		iter.fieldIdStack = iter.fieldIdStack[:len(iter.fieldIdStack)-1]
		return protocol.TType(firstByte), 0
	}
	// mask off the 4 MSB of the type header. it could contain a field id delta.
	modifier := int16((firstByte & 0xf0) >> 4)
	var fieldId protocol.FieldId
	if modifier == 0 {
		// not a delta. look ahead for the zigzag varint field id.
		fieldId = protocol.FieldId(iter.ReadInt16())
	} else {
		// has a delta. add the delta to the last read field id.
		fieldId = iter.lastFieldId + protocol.FieldId(modifier)
	}
	var fieldType protocol.TType
	if TCompactType(firstByte&0x0f) == TypeBooleanTrue {
		fieldType = protocol.BOOL
		iter.pendingBoolField = 1
	} else if TCompactType(firstByte&0x0f) == TypeBooleanFalse {
		fieldType = protocol.BOOL
		iter.pendingBoolField = 2
	} else {
		fieldType = TCompactType(firstByte & 0x0f).ToTType()
		iter.pendingBoolField = 0
	}

	// push the new field onto the field stack so we can keep the deltas going.
	iter.lastFieldId = fieldId
	return fieldType, fieldId
}

func (iter *Iterator) ReadListHeader() (protocol.TType, int) {
	lenAndType := iter.buf[0]
	iter.consume(1)
	length := int((lenAndType >> 4) & 0x0f)
	if length == 15 {
		length2 := iter.readVarInt32()
		if length2 < 0 {
			iter.ReportError("ReadListHeader", "invalid data length")
			return protocol.STOP, 0
		}
		length = int(length2)
	}
	elemType := TCompactType(lenAndType).ToTType()
	return elemType, length
}

func (iter *Iterator) ReadMapHeader() (protocol.TType, protocol.TType, int) {
	length := int(iter.readVarInt32())
	if length == 0 {
		return protocol.STOP, protocol.STOP, length
	}
	keyAndElemType := iter.buf[0]
	iter.consume(1)
	keyType := TCompactType(keyAndElemType >> 4).ToTType()
	elemType := TCompactType(keyAndElemType & 0xf).ToTType()
	return keyType, elemType, length
}

func (iter *Iterator) ReadBool() bool {
	if iter.pendingBoolField == 0 {
		return iter.ReadUInt8() == 1
	}
	return iter.pendingBoolField == 1
}

func (iter *Iterator) ReadUInt8() uint8 {
	b := iter.buf
	value := b[0]
	iter.consume(1)
	return value
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUInt8())
}

func (iter *Iterator) ReadUInt16() uint16 {
	return uint16(iter.ReadUInt32())
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadInt32())
}

func (iter *Iterator) ReadUInt32() uint32 {
	return uint32(iter.ReadInt32())
}

func (iter *Iterator) ReadInt32() int32 {
	result := iter.readVarInt32()
	u := uint32(result)
	return int32(u>>1) ^ -(result & 1)
}

func (iter *Iterator) readVarInt32() int32 {
	return int32(iter.readVarInt64())
}

func (iter *Iterator) ReadInt64() int64 {
	result := iter.readVarInt64()
	u := uint64(result)
	return int64(u>>1) ^ -(result & 1)
}

func (iter *Iterator) ReadUInt64() uint64 {
	return uint64(iter.ReadInt64())
}

func (iter *Iterator) readVarInt64() int64 {
	shift := uint(0)
	result := int64(0)
	for i, b := range iter.buf {
		result |= int64(b&0x7f) << shift
		if (b & 0x80) != 0x80 {
			iter.consume(i + 1)
			break
		}
		shift += 7
	}
	return result
}

func (iter *Iterator) ReadFloat64() float64 {
	value := math.Float64frombits(binary.LittleEndian.Uint64(iter.buf))
	iter.consume(8)
	return value
}

func (iter *Iterator) ReadString() string {
	length := iter.readVarInt32()
	value := string(iter.buf[:length])
	iter.consume(int(length))
	return value
}

func (iter *Iterator) ReadBinary() []byte {
	length := iter.readVarInt32()
	value := iter.buf[:length]
	iter.consume(int(length))
	return value
}

func (iter *Iterator) ReadStruct() map[protocol.FieldId]interface{} {
	obj := map[protocol.FieldId]interface{}{}
	iter.ReadStructHeader()
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.STOP {
			return obj
		}
		obj[fieldId] = iter.Read(fieldType)
	}
}

func (iter *Iterator) ReadList() []interface{} {
	var obj []interface{}
	elemType, length := iter.ReadListHeader()
	elemReader := iter.ReaderOf(elemType)
	for i := 0; i < length; i++ {
		obj = append(obj, elemReader())
	}
	return obj
}

func (iter *Iterator) ReadMap() map[interface{}]interface{} {
	obj := map[interface{}]interface{}{}
	keyType, elemType, length := iter.ReadMapHeader()
	keyReader := iter.ReaderOf(keyType)
	elemReader := iter.ReaderOf(elemType)
	for i := 0; i < length; i++ {
		obj[keyReader()] = elemReader()
	}
	return obj
}

func (iter *Iterator) Read(ttype protocol.TType) interface{} {
	switch ttype {
	case protocol.BOOL:
		return iter.ReadBool()
	case protocol.I08:
		return iter.ReadInt8()
	case protocol.I16:
		return iter.ReadInt16()
	case protocol.I32:
		return iter.ReadInt32()
	case protocol.I64:
		return iter.ReadInt64()
	case protocol.DOUBLE:
		return iter.ReadFloat64()
	case protocol.STRING:
		return iter.ReadString()
	case protocol.LIST:
		return iter.ReadList()
	case protocol.MAP:
		return iter.ReadMap()
	case protocol.STRUCT:
		return iter.ReadStruct()
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) ReaderOf(ttype protocol.TType) func() interface{} {
	switch ttype {
	case protocol.BOOL:
		return func() interface{} {
			return iter.ReadBool()
		}
	case protocol.I08:
		return func() interface{} {
			return iter.ReadInt8()
		}
	case protocol.I16:
		return func() interface{} {
			return iter.ReadInt16()
		}
	case protocol.I32:
		return func() interface{} {
			return iter.ReadInt32()
		}
	case protocol.I64:
		return func() interface{} {
			return iter.ReadInt64()
		}
	case protocol.DOUBLE:
		return func() interface{} {
			return iter.ReadFloat64()
		}
	case protocol.STRING:
		return func() interface{} {
			return iter.ReadString()
		}
	case protocol.LIST:
		return func() interface{} {
			return iter.ReadList()
		}
	case protocol.MAP:
		return func() interface{} {
			return iter.ReadMap()
		}
	case protocol.STRUCT:
		return func() interface{} {
			return iter.ReadStruct()
		}
	default:
		panic("unsupported type")
	}
}
