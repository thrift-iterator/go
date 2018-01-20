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

const compactProtocolId = 0x082
const compactVersion = 1
const versionMask = 0x1f

func (iter *Iterator) ReadMessageHeader() protocol.MessageHeader {
	protocolId := iter.buf[0]
	if compactProtocolId != protocolId {
		iter.ReportError("ReadMessageHeader", "invalid protocol")
		return protocol.MessageHeader{}
	}
	versionAndType := iter.buf[1]
	iter.consume(2)
	version := versionAndType & versionMask
	messageType := protocol.TMessageType((versionAndType >> 5) & 0x07)
	if version != compactVersion {
		iter.ReportError("ReadMessageHeader", fmt.Sprintf("Expected version %02x but got %02x", compactVersion, version))
		return protocol.MessageHeader{}
	}
	seqId := protocol.SeqId(iter.readVarInt32())
	messageName := iter.ReadString()
	return protocol.MessageHeader{
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
		if fieldType == protocol.TypeStop {
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
		fieldType = protocol.TypeBool
		iter.pendingBoolField = 1
	} else if TCompactType(firstByte&0x0f) == TypeBooleanFalse {
		fieldType = protocol.TypeBool
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
			return protocol.TypeStop, 0
		}
		length = int(length2)
	}
	elemType := TCompactType(lenAndType).ToTType()
	return elemType, length
}

func (iter *Iterator) ReadMapHeader() (protocol.TType, protocol.TType, int) {
	length := int(iter.readVarInt32())
	if length == 0 {
		return protocol.TypeStop, protocol.TypeStop, length
	}
	keyAndElemType := iter.buf[0]
	iter.consume(1)
	keyType := TCompactType(keyAndElemType >> 4).ToTType()
	elemType := TCompactType(keyAndElemType & 0xf).ToTType()
	return keyType, elemType, length
}

func (iter *Iterator) ReadBool() bool {
	if iter.pendingBoolField == 0 {
		return iter.ReadUint8() == 1
	}
	return iter.pendingBoolField == 1
}

func (iter *Iterator) ReadUint8() uint8 {
	b := iter.buf
	value := b[0]
	iter.consume(1)
	return value
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUint8())
}

func (iter *Iterator) ReadUint16() uint16 {
	return uint16(iter.ReadUint32())
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadInt32())
}

func (iter *Iterator) ReadUint32() uint32 {
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

func (iter *Iterator) ReadUint64() uint64 {
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

func (iter *Iterator) ReadInt() int {
	return int(iter.ReadInt64())
}

func (iter *Iterator) ReadUint() uint {
	return uint(iter.ReadUint64())
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
		if fieldType == protocol.TypeStop {
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
	case protocol.TypeBool:
		return iter.ReadBool()
	case protocol.TypeI08:
		return iter.ReadInt8()
	case protocol.TypeI16:
		return iter.ReadInt16()
	case protocol.TypeI32:
		return iter.ReadInt32()
	case protocol.TypeI64:
		return iter.ReadInt64()
	case protocol.TypeDouble:
		return iter.ReadFloat64()
	case protocol.TypeString:
		return iter.ReadString()
	case protocol.TypeList:
		return iter.ReadList()
	case protocol.TypeMap:
		return iter.ReadMap()
	case protocol.TypeStruct:
		return iter.ReadStruct()
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) ReaderOf(ttype protocol.TType) func() interface{} {
	switch ttype {
	case protocol.TypeBool:
		return func() interface{} {
			return iter.ReadBool()
		}
	case protocol.TypeI08:
		return func() interface{} {
			return iter.ReadInt8()
		}
	case protocol.TypeI16:
		return func() interface{} {
			return iter.ReadInt16()
		}
	case protocol.TypeI32:
		return func() interface{} {
			return iter.ReadInt32()
		}
	case protocol.TypeI64:
		return func() interface{} {
			return iter.ReadInt64()
		}
	case protocol.TypeDouble:
		return func() interface{} {
			return iter.ReadFloat64()
		}
	case protocol.TypeString:
		return func() interface{} {
			return iter.ReadString()
		}
	case protocol.TypeList:
		return func() interface{} {
			return iter.ReadList()
		}
	case protocol.TypeMap:
		return func() interface{} {
			return iter.ReadMap()
		}
	case protocol.TypeStruct:
		return func() interface{} {
			return iter.ReadStruct()
		}
	default:
		panic("unsupported type")
	}
}
