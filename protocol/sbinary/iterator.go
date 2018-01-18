package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/protocol/binary"
	"fmt"
)

type Iterator struct {
	real   *binary.Iterator
	reader io.Reader
	tmp    []byte
	space  []byte
	err    error
}

func NewIterator(reader io.Reader, buf []byte) *Iterator {
	return &Iterator{
		reader: reader, real: binary.NewIterator(nil),
		tmp:    make([]byte, 10),
		space:  buf,
	}
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

func (iter *Iterator) Error() error {
	return iter.err
}

func (iter *Iterator) ReportError(operation string, err string) {
	if iter.err == nil {
		iter.err = fmt.Errorf("%s: %s", operation, err)
	}
}

func (iter *Iterator) ReadMessageHeader() protocol.MessageHeader {
	versionAndMessageType := iter.ReadInt32()
	messageType := protocol.TMessageType(versionAndMessageType & 0x0ff)
	version := int64(int64(versionAndMessageType) & 0xffff0000)
	if version != protocol.VERSION_1 {
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

func (iter *Iterator) ReadMessage() protocol.Message {
	buf := iter.SkipMessage(iter.space[:0])
	if iter.err != nil {
		return protocol.Message{}
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadMessage()
}

func (iter *Iterator) ReadStructCB(cb func(fieldType protocol.TType, fieldId protocol.FieldId)) {
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.STOP {
			return
		}
		cb(fieldType, fieldId)
	}
}

func (iter *Iterator) ReadStructHeader() {
	// noop
}

func (iter *Iterator) ReadStructField() (protocol.TType, protocol.FieldId) {
	tmp := iter.tmp[:1]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadStructField", err.Error())
		return protocol.STOP, 0
	}
	fieldType := protocol.TType(tmp[0])
	if fieldType == protocol.STOP {
		return protocol.STOP, 0
	}
	fieldId := protocol.FieldId(iter.ReadUint16())
	return fieldType, fieldId
}

func (iter *Iterator) ReadStruct() map[protocol.FieldId]interface{} {
	buf := iter.SkipStruct(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadStruct()
}

func (iter *Iterator) ReadListHeader() (elemType protocol.TType, length int) {
	tmp := iter.tmp[:5]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadListHeader", err.Error())
		return protocol.STOP, 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadListHeader()
}

func (iter *Iterator) ReadList() []interface{} {
	buf := iter.SkipList(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadList()
}
func (iter *Iterator) ReadMapHeader() (keyType protocol.TType, elemType protocol.TType, length int) {
	tmp := iter.tmp[:6]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadMapHeader", err.Error())
		return protocol.STOP, protocol.STOP, 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadMapHeader()
}

func (iter *Iterator) ReadMap() map[interface{}]interface{} {
	buf := iter.SkipMap(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadMap()
}

func (iter *Iterator) ReadBool() bool {
	tmp := iter.tmp[:1]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBool", err.Error())
		return false
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadBool()
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUint8())
}

func (iter *Iterator) ReadUint8() uint8 {
	tmp := iter.tmp[:1]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUint8", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint8()
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadUint16())
}

func (iter *Iterator) ReadUint16() uint16 {
	tmp := iter.tmp[:2]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUint16", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint16()
}

func (iter *Iterator) ReadInt32() int32 {
	return int32(iter.ReadUint32())
}

func (iter *Iterator) ReadUint32() uint32 {
	tmp := iter.tmp[:4]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUint32", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint32()
}

func (iter *Iterator) ReadInt64() int64 {
	tmp := iter.allocate(8)
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadInt64", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadInt64()
}

func (iter *Iterator) ReadUint64() uint64 {
	tmp := iter.tmp[:8]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUint64", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUint64()
}

func (iter *Iterator) ReadFloat64() float64 {
	tmp := iter.tmp[:8]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadFloat64", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadFloat64()
}

func (iter *Iterator) ReadString() string {
	size := iter.ReadUint32()
	tmp := iter.allocate(int(size))
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBinary", err.Error())
		return ""
	}
	return string(tmp)
}

func (iter *Iterator) ReadBinary() []byte {
	size := iter.ReadUint32()
	tmp := iter.allocate(int(size))
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBinary", err.Error())
		return nil
	}
	return tmp
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
